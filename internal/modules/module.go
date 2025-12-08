package modules

import (
	"github.com/funvibe/funxy/internal/ast"
	"github.com/funvibe/funxy/internal/symbols"
	"github.com/funvibe/funxy/internal/typesystem"
)

// Module represents a loaded package consisting of multiple source files.
type Module struct {
	Name        string
	Dir         string
	Files       []*ast.Program
	SymbolTable *symbols.SymbolTable
	Exports     map[string]bool              // Set of exported symbol names (local + resolved re-exports)
	Imports     map[string]*Module           // Map alias/name -> Module
	TypeMap     map[ast.Node]typesystem.Type // Type inference results
	IsVirtual   bool                         // True if this is a virtual (built-in) package

	// Re-export specifications from package declaration
	// Stored during loading, resolved during analysis
	ReexportSpecs []*ast.ExportSpec

	// Evaluator-specific data
	ClassImplementations map[string]map[string]interface{} // Runtime trait implementations

	// Package group support (import "dir" imports all dir/*)
	IsPackageGroup bool     // True if this is a combined package from subdirectories
	SubPackages    []string // Names of sub-packages (e.g., ["utils", "helpers"])

	HeadersAnalyzed  bool
	HeadersAnalyzing bool
	BodiesAnalyzed   bool
	BodiesAnalyzing  bool
}

// GetExports returns a map of exported symbol names to their runtime values.
// Implements evaluator.LoadedModule interface.
// BUT Module doesn't store runtime values (evaluator.Object).
// Values are computed when module is evaluated.
// Module struct only holds AST and SymbolTable.
// We need a way to get EVALUATED exports.
// Loader.LoadedModules stores *Module.
// *Module doesn't have values.
// We need to EVALUATE the module if not already evaluated.
// But Evaluator shouldn't evaluate modules recursively here if Loader already handles caching?
// Loader just loads AST.
// Main.go `evaluateModule` evaluates it and caches the result in `moduleCache`.
// BUT `moduleCache` is local to `main.go`.
// Evaluator needs access to evaluated modules.
// `Loader` interface returns `interface{}`.
// If `Evaluator` uses `Loader` to get module, it gets `*Module` (AST).
// It needs to EVALUATE it.
// So `evalImportStatement` must call `evaluateModule` logic.
// BUT `evaluateModule` logic is in `main.go`.
// We should move `evaluateModule` logic to `Evaluator` or `Loader`.
// Or `Loader` should return evaluated module? No, Loader loads AST.
//
// Best place: `Evaluator` should have a method to recursively evaluate modules.
// `Evaluator` already has `Eval`.
// We can implement `EvaluateModule` inside `Evaluator`.
// And we need a cache for evaluated modules in `Evaluator` instance?
// Yes.

// GetExportedSymbols returns a map of exported symbol names to their types.
// Implements analyzer.LoadedModule interface.
func (m *Module) GetExportedSymbols() map[string]typesystem.Type {
	result := make(map[string]typesystem.Type)
	for name := range m.Exports {
		if sym, ok := m.SymbolTable.Find(name); ok {
			result[name] = sym.Type
		}
	}
	return result
}

func (m *Module) GetExports() map[string]symbols.Symbol {
	result := make(map[string]symbols.Symbol)
	
	// For package groups, collect exports from all sub-packages
	if m.IsPackageGroup {
		for _, subMod := range m.Imports {
			for name := range subMod.Exports {
				if sym, ok := subMod.SymbolTable.Find(name); ok {
					result[name] = sym
				}
			}
		}
		return result
	}
	
	// Regular module: exports are resolved (local + re-exports added by analyzer)
	for name := range m.Exports {
		if sym, ok := m.SymbolTable.Find(name); ok {
			result[name] = sym
		}
	}
	
	return result
}

// GetReexportSpecs returns the re-export specifications from package declaration.
func (m *Module) GetReexportSpecs() []*ast.ExportSpec {
	return m.ReexportSpecs
}

// AddExport adds a symbol name to the exports set.
// Used by analyzer to add resolved re-exports.
func (m *Module) AddExport(name string) {
	if m.Exports == nil {
		m.Exports = make(map[string]bool)
	}
	m.Exports[name] = true
}

func (m *Module) GetName() string             { return m.Name }
func (m *Module) IsHeadersAnalyzed() bool    { return m.HeadersAnalyzed }
func (m *Module) SetHeadersAnalyzed(v bool)  { m.HeadersAnalyzed = v }
func (m *Module) IsHeadersAnalyzing() bool   { return m.HeadersAnalyzing }
func (m *Module) SetHeadersAnalyzing(v bool) { m.HeadersAnalyzing = v }

func (m *Module) IsBodiesAnalyzed() bool    { return m.BodiesAnalyzed }
func (m *Module) SetBodiesAnalyzed(v bool)  { m.BodiesAnalyzed = v }
func (m *Module) IsBodiesAnalyzing() bool   { return m.BodiesAnalyzing }
func (m *Module) SetBodiesAnalyzing(v bool) { m.BodiesAnalyzing = v }

func (m *Module) IsAnalyzed() bool {
	return m.BodiesAnalyzed
}

func (m *Module) SetAnalyzed(v bool) {
	m.BodiesAnalyzed = v
}

func (m *Module) IsAnalyzing() bool {
	return m.BodiesAnalyzing
}

func (m *Module) SetAnalyzing(v bool) {
	m.BodiesAnalyzing = v
}

// IsPackageGroupModule returns true if this module is a combined package from subdirectories
func (m *Module) IsPackageGroupModule() bool {
	return m.IsPackageGroup
}

// GetSubModulesRaw returns sub-modules for package groups (implements LoadedModule interface)
// Note: returns interface{} because we can't import analyzer.LoadedModule here
func (m *Module) GetSubModulesRaw() map[string]interface{} {
	result := make(map[string]interface{})
	for name, mod := range m.Imports {
		result[name] = mod
	}
	return result
}

func (m *Module) GetFiles() []*ast.Program {
	return m.Files
}

func (m *Module) GetSymbolTable() *symbols.SymbolTable {
	return m.SymbolTable
}

func (m *Module) SetTypeMap(tm map[ast.Node]typesystem.Type) {
	m.TypeMap = tm
}

// NewModule creates a new empty module.
func NewModule(name, dir string) *Module {
	return &Module{
		Name:               name,
		Dir:                dir,
		Files:              []*ast.Program{},
		SymbolTable:        symbols.NewSymbolTable(), // Each module has its own scope
		Exports:            make(map[string]bool),
		Imports:            make(map[string]*Module),
		ClassImplementations: make(map[string]map[string]interface{}),
	}
}
