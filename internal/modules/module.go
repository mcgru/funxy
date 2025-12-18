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

	// Trait default implementations found during analysis
	TraitDefaults map[string]*ast.FunctionStatement

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
func (m *Module) GetExports() map[string]symbols.Symbol {
	exportedSymbols := make(map[string]symbols.Symbol)

	if m.IsPackageGroup {
		// For package groups, look up symbols in sub-modules
		for name := range m.Exports {
			// Find which sub-module exports this symbol
			for _, subMod := range m.Imports {
				if subMod.Exports[name] {
					if sym, ok := subMod.SymbolTable.Find(name); ok {
						exportedSymbols[name] = sym
						break // Found it
					}
				}
			}
		}
	} else {
		// Regular module
		for name := range m.Exports {
			if sym, ok := m.SymbolTable.Find(name); ok {
				exportedSymbols[name] = sym
			}
		}
	}
	return exportedSymbols
}

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) GetExportedSymbols() map[string]typesystem.Type {
	return nil // Deprecated
}

func (m *Module) IsHeadersAnalyzed() bool {
	return m.HeadersAnalyzed
}

func (m *Module) SetHeadersAnalyzed(v bool) {
	m.HeadersAnalyzed = v
}

func (m *Module) IsHeadersAnalyzing() bool {
	return m.HeadersAnalyzing
}

func (m *Module) SetHeadersAnalyzing(v bool) {
	m.HeadersAnalyzing = v
}

func (m *Module) IsBodiesAnalyzed() bool {
	return m.BodiesAnalyzed
}

func (m *Module) SetBodiesAnalyzed(v bool) {
	m.BodiesAnalyzed = v
}

func (m *Module) IsBodiesAnalyzing() bool {
	return m.BodiesAnalyzing
}

func (m *Module) SetBodiesAnalyzing(v bool) {
	m.BodiesAnalyzing = v
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

func (m *Module) IsPackageGroupModule() bool {
	return m.IsPackageGroup
}

func (m *Module) GetSubModulesRaw() map[string]interface{} {
	subs := make(map[string]interface{})
	for k, v := range m.Imports {
		subs[k] = v
	}
	return subs
}

func (m *Module) GetReexportSpecs() []*ast.ExportSpec {
	return m.ReexportSpecs
}

func (m *Module) AddExport(name string) {
	m.Exports[name] = true
}

func (m *Module) SetTraitDefaults(defaults map[string]*ast.FunctionStatement) {
	m.TraitDefaults = defaults
}

func (m *Module) GetTraitDefaults() map[string]*ast.FunctionStatement {
	return m.TraitDefaults
}
