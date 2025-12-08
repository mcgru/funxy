package analyzer

import (
	"fmt"
	"github.com/funvibe/funxy/internal/ast"
	"github.com/funvibe/funxy/internal/config"
	"github.com/funvibe/funxy/internal/diagnostics"
	"github.com/funvibe/funxy/internal/symbols"
	"github.com/funvibe/funxy/internal/token"
	"github.com/funvibe/funxy/internal/typesystem"
	"github.com/funvibe/funxy/internal/utils"
	"sort"
)

func (w *walker) VisitPackageDeclaration(n *ast.PackageDeclaration) {}

func tagModule(t typesystem.Type, moduleName string, exportedTypes map[string]bool) typesystem.Type {
	if t == nil {
		return nil
	}

	switch t := t.(type) {
	case typesystem.TCon:
		if exportedTypes[t.Name] {
			t.Module = moduleName
		}
		return t
	case typesystem.TApp:
		newConstructor := tagModule(t.Constructor, moduleName, exportedTypes)
		newArgs := []typesystem.Type{}
		for _, arg := range t.Args {
			newArgs = append(newArgs, tagModule(arg, moduleName, exportedTypes))
		}
		return typesystem.TApp{Constructor: newConstructor, Args: newArgs}
	case typesystem.TFunc:
		newParams := []typesystem.Type{}
		for _, p := range t.Params {
			newParams = append(newParams, tagModule(p, moduleName, exportedTypes))
		}
		newRet := tagModule(t.ReturnType, moduleName, exportedTypes)
		return typesystem.TFunc{
			Params:       newParams,
			ReturnType:   newRet,
			IsVariadic:   t.IsVariadic,
			DefaultCount: t.DefaultCount,
			Constraints:  t.Constraints,
		}
	case typesystem.TTuple:
		newElems := []typesystem.Type{}
		for _, el := range t.Elements {
			newElems = append(newElems, tagModule(el, moduleName, exportedTypes))
		}
		return typesystem.TTuple{Elements: newElems}
	case typesystem.TRecord:
		newFields := make(map[string]typesystem.Type)
		for k, v := range t.Fields {
			newFields[k] = tagModule(v, moduleName, exportedTypes)
		}
		return typesystem.TRecord{Fields: newFields, IsOpen: t.IsOpen}
	case typesystem.TType:
		return typesystem.TType{Type: tagModule(t.Type, moduleName, exportedTypes)}
	}
	return t
}

func (w *walker) VisitImportStatement(n *ast.ImportStatement) {
	if w.loader == nil {
		return
	}

	// Resolve absolute path from ImportStatement
	importPath := n.Path.Value
	pathToCheck := utils.ResolveImportPath(w.BaseDir, importPath)

	modInterface, err := w.loader.GetModule(pathToCheck)
	if err != nil {
		// Check if it's a DiagnosticError (syntax error in module)
		if compileErr, ok := err.(*diagnostics.DiagnosticError); ok {
			w.addError(compileErr)
			return
		}

		w.addError(diagnostics.NewError(
			diagnostics.ErrA001,
			n.Path.GetToken(),
			"module not found: "+n.Path.Value,
		))
		return
	}

	// Define Module Symbol
	name := ""
	if n.Alias != nil {
		name = n.Alias.Value
	} else {
		name = utils.ExtractModuleName(n.Path.Value)
	}

	if loadedMod, ok := modInterface.(LoadedModule); ok {
		// Handle package groups by analyzing all sub-packages
		if loadedMod.IsPackageGroupModule() {
			w.analyzePackageGroup(loadedMod, pathToCheck)
		} else {
			// Recursive Analysis based on Mode for regular modules
			w.analyzeRegularModule(loadedMod, pathToCheck)
		}
		
		// Store mapping alias -> packageName for extension method/trait lookup
		packageName := loadedMod.GetName()
		w.symbolTable.RegisterModuleAlias(name, packageName)
		
		// Get Exports (Symbols with Kinds)
		exportSymbols := loadedMod.GetExports()
		
		// Identify exported Types to tag TCon (sorted for deterministic order)
		exportedTypes := make(map[string]bool)
		exportKeys := make([]string, 0, len(exportSymbols))
		for k := range exportSymbols {
			exportKeys = append(exportKeys, k)
		}
		sort.Strings(exportKeys)
		
		for _, expName := range exportKeys {
			sym := exportSymbols[expName]
			if sym.Kind == symbols.TypeSymbol || sym.Kind == symbols.ConstructorSymbol { 
				// Only tag TypeSymbol? ConstructorSymbol usually not TCon (TFunc or TApp)
				// But if we export Type, we want TCon{Name: Type} to be tagged.
				if sym.Kind == symbols.TypeSymbol {
					exportedTypes[expName] = true
				}
			}
		}

		// Handle selective imports
		if n.ImportAll || len(n.Symbols) > 0 || len(n.Exclude) > 0 {
			// Import symbols directly into current scope
			symbolsToImport := make(map[string]bool)
			
			if n.ImportAll {
				// Import all (using already sorted exportKeys)
				for _, expName := range exportKeys {
					symbolsToImport[expName] = true
				}
			} else if len(n.Symbols) > 0 {
				// Import only specified
				for _, sym := range n.Symbols {
					if _, ok := exportSymbols[sym.Value]; ok {
						symbolsToImport[sym.Value] = true
					} else {
						w.addError(diagnostics.NewError(
							diagnostics.ErrA006,
							sym.GetToken(),
							sym.Value,
						))
					}
				}
			} else if len(n.Exclude) > 0 {
				// Import all except specified (using already sorted exportKeys)
				excludeSet := make(map[string]bool)
				for _, sym := range n.Exclude {
					excludeSet[sym.Value] = true
				}
				for _, expName := range exportKeys {
					if !excludeSet[expName] {
						symbolsToImport[expName] = true
					}
				}
			}
			
			// Register each symbol directly (sorted for deterministic order)
			importKeys := make([]string, 0, len(symbolsToImport))
			for k := range symbolsToImport {
				importKeys = append(importKeys, k)
			}
			sort.Strings(importKeys)

		for _, symName := range importKeys {
			sym := exportSymbols[symName]
			taggedType := tagModule(sym.Type, packageName, exportedTypes)
			
			// Use OriginModule from source symbol, or packageName if not set
			origin := sym.OriginModule
			if origin == "" {
				origin = packageName
			}
			
			// Check for duplicate import - allow if same origin (same symbol re-exported via different paths)
			if existing, exists := w.symbolTable.Find(symName); exists {
				if existing.OriginModule == origin {
					// Same symbol from same origin - OK, skip
					continue
				}
				// Different origins - conflict
				w.addError(diagnostics.NewError(
					diagnostics.ErrA004,
					n.GetToken(),
					fmt.Sprintf("%s' (already defined from %s, cannot import from %s)", 
						symName, existing.OriginModule, origin),
				))
				continue
			}
			
			if sym.Kind == symbols.TypeSymbol {
				// Check if it's a type alias with underlying type
				if sym.IsTypeAlias() {
					// Copy both nominal type and underlying type
					w.symbolTable.DefineTypeAlias(symName, taggedType, sym.UnderlyingType, origin)
				} else {
					w.symbolTable.DefineType(symName, taggedType, origin)
				}
			} else if sym.Kind == symbols.ConstructorSymbol {
				w.symbolTable.DefineConstructor(symName, taggedType, origin)
			} else {
				w.symbolTable.Define(symName, taggedType, origin)
			}
		}
		
		// Copy trait implementations and extension methods for imported types
		if loadedMod, ok := modInterface.(LoadedModule); ok {
			if modSymTable := loadedMod.GetSymbolTable(); modSymTable != nil {
				importedTypes := make(map[string]bool)
				for _, symName := range importKeys {
					if exportedTypes[symName] {
						importedTypes[symName] = true
					}
				}
				w.importTraitImplementations(modSymTable, importedTypes, "")
				w.importExtensionMethods(modSymTable, importedTypes)
			}
		}
		} else {
			// Default: Create TRecord for the module (using already sorted exportKeys)
			fields := make(map[string]typesystem.Type)
			for _, expName := range exportKeys {
				sym := exportSymbols[expName]
				if sym.Type == nil {
					// Skip symbols with nil types (e.g., traits)
					continue
				}
				// Tag with alias for namespaced access (e.g., m.Vector)
				taggedType := tagModule(sym.Type, name, exportedTypes)
				if taggedType == nil {
					continue
				}
				fields[expName] = taggedType
			}

			moduleType := typesystem.TRecord{Fields: fields}
			w.symbolTable.DefineModule(name, moduleType)

			// Copy type aliases for exported types (qualified names like m.Vector)
			if modSymTable := loadedMod.GetSymbolTable(); modSymTable != nil {
				for typeName := range exportedTypes {
					if aliasType, ok := modSymTable.GetTypeAlias(typeName); ok {
						// Register with qualified name (e.g., "m.Vector" -> TRecord{...})
						w.symbolTable.RegisterTypeAlias(name+"."+typeName, aliasType)
					}
				}
			}

			// Copy trait implementations for exported types
			// This is still needed for types that don't have Module field (like record aliases)
			// For types with Module field, we also look up in source module via isImplementationInSourceModule
			if modSymTable := loadedMod.GetSymbolTable(); modSymTable != nil {
				w.importTraitImplementations(modSymTable, exportedTypes, name)
			}
		}
	} else {
		// Fallback if casting fails (shouldn't happen if wired correctly)
		w.symbolTable.Define(name, typesystem.TCon{Name: "Module"}, name)
	}
}

// analyzePackageGroup analyzes all sub-packages in a package group
func (w *walker) analyzePackageGroup(loadedMod LoadedModule, pathToCheck string) {
	// Analyze each sub-package
	for _, subModRaw := range loadedMod.GetSubModulesRaw() {
		subMod, ok := subModRaw.(LoadedModule)
		if !ok {
			continue
		}
		w.analyzeRegularModule(subMod, pathToCheck)
	}
	
	// Mark the package group as analyzed
	loadedMod.SetHeadersAnalyzed(true)
	loadedMod.SetBodiesAnalyzed(true)
}

// resolveReexports resolves re-export specifications and adds symbols to module exports.
// Called after module analysis when all imports have been processed.
// Returns errors if re-export references a module that wasn't imported.
func resolveReexports(mod LoadedModule) []*diagnostics.DiagnosticError {
	var errors []*diagnostics.DiagnosticError
	
	specs := mod.GetReexportSpecs()
	if len(specs) == 0 {
		return errors
	}
	
	symbolTable := mod.GetSymbolTable()
	
	for _, spec := range specs {
		if spec.ModuleName == nil {
			continue
		}
		
		moduleName := spec.ModuleName.Value
		
		// Validate that the module was actually imported
		// Check if moduleName is a registered module alias
		_, isImported := symbolTable.GetPackageNameByAlias(moduleName)
		if !isImported {
			errors = append(errors, diagnostics.NewError(
				diagnostics.ErrA006,
				spec.GetToken(),
				fmt.Sprintf("cannot re-export from '%s': module not imported", moduleName),
			))
			continue
		}
		
		if spec.ReexportAll {
			// shapes(*) — re-export all symbols with OriginModule == moduleName
			// We iterate through all symbols and find those that came from this module
			for name, sym := range symbolTable.All() {
				if sym.OriginModule == moduleName {
					mod.AddExport(name)
				}
			}
		} else {
			// shapes(foo, bar) — re-export specific symbols
			for _, symIdent := range spec.Symbols {
				name := symIdent.Value
				// Verify symbol exists in symbol table
				if _, ok := symbolTable.Find(name); ok {
					mod.AddExport(name)
				}
			}
		}
	}
	
	return errors
}

// analyzeRegularModule handles analysis for a regular (non-package-group) module
func (w *walker) analyzeRegularModule(loadedMod LoadedModule, pathToCheck string) {
	if w.mode == ModeHeaders {
		if !loadedMod.IsHeadersAnalyzed() {
			if loadedMod.IsHeadersAnalyzing() {
				// Circular dependency in headers - skip
				return
			}
			loadedMod.SetHeadersAnalyzing(true)

			modAnalyzer := New(loadedMod.GetSymbolTable())
			modAnalyzer.SetLoader(w.loader)
			modAnalyzer.RegisterBuiltins()
			modAnalyzer.BaseDir = utils.GetModuleDir(pathToCheck)

			for _, file := range loadedMod.GetFiles() {
				errs := modAnalyzer.AnalyzeHeaders(file)
				w.addErrors(errs)
			}

			loadedMod.SetTypeMap(modAnalyzer.TypeMap)
			loadedMod.SetHeadersAnalyzing(false)
			loadedMod.SetHeadersAnalyzed(true)
			
			// Resolve re-exports after headers analysis (imports are processed)
			reexportErrs := resolveReexports(loadedMod)
			w.addErrors(reexportErrs)
		}
	} else if w.mode == ModeBodies {
		if !loadedMod.IsBodiesAnalyzed() {
			if loadedMod.IsBodiesAnalyzing() {
				// Cycle in bodies is allowed - skip
				return
			}
			loadedMod.SetBodiesAnalyzing(true)

			modAnalyzer := New(loadedMod.GetSymbolTable())
			modAnalyzer.SetLoader(w.loader)
			modAnalyzer.RegisterBuiltins()
			modAnalyzer.BaseDir = utils.GetModuleDir(pathToCheck)

			for _, file := range loadedMod.GetFiles() {
				errs := modAnalyzer.AnalyzeBodies(file)
				w.addErrors(errs)
			}

			loadedMod.SetTypeMap(modAnalyzer.TypeMap)
			loadedMod.SetBodiesAnalyzing(false)
			loadedMod.SetBodiesAnalyzed(true)
		}
	} else {
		// Legacy / Full Mode
		if !loadedMod.IsBodiesAnalyzed() {
			if loadedMod.IsBodiesAnalyzing() {
				return
			}
			loadedMod.SetBodiesAnalyzing(true)

			modAnalyzer := New(loadedMod.GetSymbolTable())
			modAnalyzer.SetLoader(w.loader)
			modAnalyzer.RegisterBuiltins()
			modAnalyzer.BaseDir = utils.GetModuleDir(pathToCheck)

			for _, file := range loadedMod.GetFiles() {
				errs := modAnalyzer.Analyze(file)
				w.addErrors(errs)
			}

			loadedMod.SetTypeMap(modAnalyzer.TypeMap)
			loadedMod.SetBodiesAnalyzing(false)
			loadedMod.SetBodiesAnalyzed(true)
			
			// Resolve re-exports after full analysis
			reexportErrs := resolveReexports(loadedMod)
			w.addErrors(reexportErrs)
		}
	}
}

func RegisterFunctionDeclaration(n *ast.FunctionStatement, table *symbols.SymbolTable, freshVar func() string, origin string) []*diagnostics.DiagnosticError {
	var errors []*diagnostics.DiagnosticError

	// Skip if function was not properly parsed
	if n == nil || n.Name == nil {
		return errors
	}

	// Check for shadowing, but allow overwriting Pending symbols (forward declarations)
	if table.IsDefined(n.Name.Value) {
		sym, ok := table.Find(n.Name.Value)
		if ok && !sym.IsPending {
			errors = append(errors, diagnostics.NewError(diagnostics.ErrA004, n.Name.GetToken(), n.Name.Value))
			return errors
		}
	}

	// 1. Create temporary scope for signature analysis (Type Params)
	sigScope := symbols.NewEnclosedSymbolTable(table)

	// 2. Register Type Params
	for _, tp := range n.TypeParams {
		sigScope.DefineType(tp.Value, typesystem.TVar{Name: tp.Value}, "")
	}

	// 3. Build Return Type
	var retType typesystem.Type
	if n.ReturnType != nil {
		retType = BuildType(n.ReturnType, sigScope, &errors)
	} else {
		retType = typesystem.TVar{Name: freshVar()}
	}

	// 4. Build Params
	var params []typesystem.Type
	// If extension, add receiver
	if n.Receiver != nil && n.Receiver.Type != nil {
		params = append(params, BuildType(n.Receiver.Type, sigScope, &errors))
	}

	var isVariadic bool
	var defaultCount int
	for _, p := range n.Parameters {
		if p.IsVariadic {
			isVariadic = true
		}
		if p.Default != nil {
			defaultCount++
		}
		if p.Type != nil {
			params = append(params, BuildType(p.Type, sigScope, &errors))
		} else {
			params = append(params, typesystem.TVar{Name: freshVar()})
		}
	}

	// Build constraints for TFunc
	var fnConstraints []typesystem.Constraint
	for _, c := range n.Constraints {
		fnConstraints = append(fnConstraints, typesystem.Constraint{TypeVar: c.TypeVar, Trait: c.Trait})
	}

	fnType := typesystem.TFunc{
		Params:       params,
		ReturnType:   retType,
		IsVariadic:   isVariadic,
		DefaultCount: defaultCount,
		Constraints:  fnConstraints,
	}

	// 5. Define in Table (Outer)
	if n.Receiver != nil {
		typeName := resolveReceiverTypeName(n.Receiver.Type, table)
		if typeName == "" {
			errors = append(errors, diagnostics.NewError(diagnostics.ErrA003, n.Receiver.Token, "invalid receiver type"))
		} else {
			table.RegisterExtensionMethod(typeName, n.Name.Value, fnType)
		}
	} else {
		table.Define(n.Name.Value, fnType, origin)
	}

	return errors
}

func (w *walker) VisitFunctionStatement(n *ast.FunctionStatement) {
	// Skip if function was not properly parsed
	if n == nil || n.Name == nil {
		return
	}
	
	// 0. Check naming convention (function name must start with lowercase)
	if n.Receiver == nil { // Only for regular functions, not extension methods
		if !checkValueName(n.Name.Value, n.Name.Token, &w.errors) {
			return
		}
	}
	
	// 1. Register Signature in Outer Scope
	
	// 1. Prepare types
	var retType typesystem.Type
	if n.ReturnType != nil {
		retType = BuildType(n.ReturnType, w.symbolTable, &w.errors)
	} else {
		retType = w.freshVar()
	}

	// 2. Create new scope for function body (and header analysis)
	outer := w.symbolTable
	w.symbolTable = symbols.NewEnclosedSymbolTable(outer)
	defer func() { w.symbolTable = outer }()

	// 4. Register Generic Constraints / Type Params
	for _, tp := range n.TypeParams {
		w.symbolTable.DefineType(tp.Value, typesystem.TVar{Name: tp.Value}, "")
	}

	var params []typesystem.Type

	// If extension, add receiver to params first
	if n.Receiver != nil && n.Receiver.Type != nil {
		params = append(params, BuildType(n.Receiver.Type, w.symbolTable, &w.errors))
	}

	var isVariadic bool
	var defaultCount int
	for _, p := range n.Parameters {
		if p.IsVariadic {
			isVariadic = true
		}
		if p.Default != nil {
			defaultCount++
		}
		if p.Type != nil {
			t := BuildType(p.Type, w.symbolTable, &w.errors)
			params = append(params, t)
		} else {
			tv := w.freshVar()
			params = append(params, tv)
		}
	}

	// Build constraints for TFunc
	var fnConstraints []typesystem.Constraint
	for _, c := range n.Constraints {
		fnConstraints = append(fnConstraints, typesystem.Constraint{TypeVar: c.TypeVar, Trait: c.Trait})
	}

	fnType := typesystem.TFunc{
		Params:       params,
		ReturnType:   retType,
		IsVariadic:   isVariadic,
		DefaultCount: defaultCount,
		Constraints:  fnConstraints,
	}

	if n.Receiver != nil {
		typeName := resolveReceiverTypeName(n.Receiver.Type, outer) 
		if typeName == "" {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				n.Receiver.Token,
				"invalid receiver type for extension method",
			))
		} else {
			outer.RegisterExtensionMethod(typeName, n.Name.Value, fnType)
		}
	} else {
		outer.Define(n.Name.Value, fnType, w.currentModuleName)
	}

	// 2.5 Register Receiver in scope
	if n.Receiver != nil {
		w.symbolTable.Define(n.Receiver.Name.Value, params[0], "")
	}

	// 3. Register parameters in the new scope
	for i, param := range n.Parameters {
		idx := i
		if n.Receiver != nil {
			idx++
		}
		
		paramType := params[idx]
		
		if param.IsVariadic {
			// Wrap in List
			paramType = typesystem.TApp{
				Constructor: typesystem.TCon{Name: config.ListTypeName},
				Args:        []typesystem.Type{paramType},
			}
		}
		w.symbolTable.Define(param.Name.Value, paramType, "")
	}

	// 5. Analyze body
	if n.Body != nil {
		prevInLoop := w.inLoop
		w.inLoop = false
		n.Body.Accept(w)
		w.markTailCalls(n.Body)
		w.inLoop = prevInLoop

		bodyType, sBody, err := InferWithContext(w.inferCtx, n.Body, w.symbolTable)
		if err != nil {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				n.Body.GetToken(),
				err.Error(),
			))
		} else {
			// Apply accumulated substitution from body to return type before unification
			retType = retType.Apply(sBody)
			
			subst, err := typesystem.Unify(retType, bodyType)
			if err != nil {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA003,
					n.Body.GetToken(),
					"function return type mismatch: declared "+retType.String()+", got "+bodyType.String(),
				))
			}
			
			// Compose all substitutions
			finalSubst := subst.Compose(sBody)
			
			// Apply to function type
			fnType = fnType.Apply(finalSubst).(typesystem.TFunc)

			// Apply substitution to all nodes in body in TypeMap
			w.applySubstToNode(n.Body, finalSubst)
		}
	}
	
	// 6. Store Function Type in TypeMap for Compiler
	w.TypeMap[n] = fnType
}

// applySubstToNode recursively applies a type substitution to all nodes in the AST.
// This ensures that type variables resolved during inference are propagated to all
// sub-expressions in the TypeMap.
func (w *walker) applySubstToNode(node ast.Node, subst typesystem.Subst) {
	if node == nil || len(subst) == 0 {
		return
	}

	// Update type in TypeMap if present
	if t, ok := w.TypeMap[node]; ok {
		w.TypeMap[node] = t.Apply(subst)
	}

	// Recursively traverse children based on node type
	switch n := node.(type) {
	// ==================== Statements ====================
	case *ast.Program:
		for _, stmt := range n.Statements {
			w.applySubstToNode(stmt, subst)
		}

	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			w.applySubstToNode(stmt, subst)
		}

	case *ast.ExpressionStatement:
		w.applySubstToNode(n.Expression, subst)

	case *ast.FunctionStatement:
		w.applySubstToNode(n.Body, subst)

	case *ast.ConstantDeclaration:
		w.applySubstToNode(n.Value, subst)

	case *ast.BreakStatement:
		if n.Value != nil {
			w.applySubstToNode(n.Value, subst)
		}

	case *ast.ContinueStatement:
		// No expression children

	case *ast.InstanceDeclaration:
		for _, method := range n.Methods {
			w.applySubstToNode(method.Body, subst)
		}

	// ==================== Expressions ====================
	// --- Operators ---
	case *ast.AssignExpression:
		w.applySubstToNode(n.Left, subst)
		w.applySubstToNode(n.Value, subst)

	case *ast.InfixExpression:
		w.applySubstToNode(n.Left, subst)
		w.applySubstToNode(n.Right, subst)

	case *ast.PrefixExpression:
		w.applySubstToNode(n.Right, subst)

	case *ast.PostfixExpression:
		w.applySubstToNode(n.Left, subst)

	// --- Function calls ---
	case *ast.CallExpression:
		w.applySubstToNode(n.Function, subst)
		for _, arg := range n.Arguments {
			w.applySubstToNode(arg, subst)
		}

	case *ast.TypeApplicationExpression:
		w.applySubstToNode(n.Expression, subst)

	case *ast.FunctionLiteral:
		w.applySubstToNode(n.Body, subst)

	// --- Control flow ---
	case *ast.IfExpression:
		w.applySubstToNode(n.Condition, subst)
		w.applySubstToNode(n.Consequence, subst)
		if n.Alternative != nil {
		w.applySubstToNode(n.Alternative, subst)
		}

	case *ast.ForExpression:
		if n.Condition != nil {
			w.applySubstToNode(n.Condition, subst)
		}
		if n.Iterable != nil {
			w.applySubstToNode(n.Iterable, subst)
		}
		if n.Body != nil {
			w.applySubstToNode(n.Body, subst)
		}

	case *ast.MatchExpression:
		w.applySubstToNode(n.Expression, subst)
		for _, arm := range n.Arms {
			w.applySubstToNode(arm.Expression, subst)
			// Also traverse patterns for any nested expressions
			w.applySubstToPattern(arm.Pattern, subst)
		}

	// --- Collection literals ---
	case *ast.TupleLiteral:
		for _, elem := range n.Elements {
			w.applySubstToNode(elem, subst)
		}

	case *ast.ListLiteral:
		for _, elem := range n.Elements {
			w.applySubstToNode(elem, subst)
		}

	case *ast.RecordLiteral:
		if n.Spread != nil {
			w.applySubstToNode(n.Spread, subst)
		}
		for _, val := range n.Fields {
			w.applySubstToNode(val, subst)
		}

	// --- Access expressions ---
	case *ast.IndexExpression:
		w.applySubstToNode(n.Left, subst)
		w.applySubstToNode(n.Index, subst)

	case *ast.MemberExpression:
		w.applySubstToNode(n.Left, subst)
		w.applySubstToNode(n.Member, subst)

	// --- Other expressions ---
	case *ast.SpreadExpression:
		w.applySubstToNode(n.Expression, subst)

	case *ast.AnnotatedExpression:
		w.applySubstToNode(n.Expression, subst)

	// --- Literals (no children, but may have TypeMap entries) ---
	case *ast.Identifier,
		*ast.IntegerLiteral,
		*ast.FloatLiteral,
		*ast.BooleanLiteral,
		*ast.NilLiteral,
		*ast.BigIntLiteral,
		*ast.RationalLiteral,
		*ast.StringLiteral,
		*ast.InterpolatedString,
		*ast.CharLiteral:
		// No children to traverse, TypeMap already updated above
	}
}

// applySubstToPattern applies substitution to patterns that may contain sub-patterns.
func (w *walker) applySubstToPattern(pattern ast.Pattern, subst typesystem.Subst) {
	if pattern == nil {
		return
	}

	// Update type in TypeMap if present
	if t, ok := w.TypeMap[pattern]; ok {
		w.TypeMap[pattern] = t.Apply(subst)
	}

	switch p := pattern.(type) {
	case *ast.ConstructorPattern:
		for _, elem := range p.Elements {
			w.applySubstToPattern(elem, subst)
		}

	case *ast.TuplePattern:
		for _, elem := range p.Elements {
			w.applySubstToPattern(elem, subst)
		}

	case *ast.ListPattern:
		for _, elem := range p.Elements {
			w.applySubstToPattern(elem, subst)
		}

	case *ast.SpreadPattern:
		w.applySubstToPattern(p.Pattern, subst)

	case *ast.RecordPattern:
		for _, pat := range p.Fields {
			w.applySubstToPattern(pat, subst)
		}

	case *ast.LiteralPattern, *ast.WildcardPattern, *ast.IdentifierPattern:
		// No children to traverse
	}
}

func RegisterTypeDeclaration(stmt *ast.TypeDeclarationStatement, table *symbols.SymbolTable, origin string) []*diagnostics.DiagnosticError {
	var errors []*diagnostics.DiagnosticError
	if stmt == nil || stmt.Name == nil {
		return errors
	}

	// Check naming convention (type name must start with uppercase)
	if !checkTypeName(stmt.Name.Value, stmt.Name.GetToken(), &errors) {
		return errors
	}

	// Check for shadowing, but allow overwriting Pending symbols (forward declarations)
	if table.IsDefined(stmt.Name.Value) {
		sym, ok := table.Find(stmt.Name.Value)
		if ok && !sym.IsPending {
			errors = append(errors, diagnostics.NewError(
				diagnostics.ErrA004,
				stmt.Name.GetToken(),
				stmt.Name.Value,
			))
			return errors
		}
		// If it is pending, we proceed to overwrite it
	}

	// Register the type name immediately
	tCon := typesystem.TCon{Name: stmt.Name.Value}
	table.DefineType(stmt.Name.Value, tCon, origin)

	// Register Kind
	kind := typesystem.Star
	if len(stmt.TypeParameters) > 0 {
		kinds := make([]typesystem.Kind, len(stmt.TypeParameters)+1)
		for i := range stmt.TypeParameters {
			kinds[i] = typesystem.Star
		}
		kinds[len(stmt.TypeParameters)] = typesystem.Star
		kind = typesystem.MakeArrow(kinds...)
	}
	table.RegisterKind(stmt.Name.Value, kind)

	// Register type parameters
	if len(stmt.TypeParameters) > 0 {
		params := make([]string, len(stmt.TypeParameters))
		for i, p := range stmt.TypeParameters {
			params[i] = p.Value
		}
		table.RegisterTypeParams(stmt.Name.Value, params)
	}

	// 1. Create temporary scope for type parameters
	typeScope := symbols.NewEnclosedSymbolTable(table)
	for _, tp := range stmt.TypeParameters {
		typeScope.DefineType(tp.Value, typesystem.TVar{Name: tp.Value}, "")
	}

	if stmt.IsAlias {
		if stmt.TargetType == nil {
			return errors
		}
		// Validate target type
		if nt, ok := stmt.TargetType.(*ast.NamedType); ok {
			if !table.IsDefined(nt.Name.Value) && !typeScope.IsDefined(nt.Name.Value) {
				errors = append(errors, diagnostics.NewError(
					diagnostics.ErrA002,
					nt.GetToken(),
					nt.Name.Value,
				))
			}
		}

		// Use typeScope to build the type
		realType := BuildType(stmt.TargetType, typeScope, &errors)
		// Use DefineTypeAlias: TCon for trait lookup, realType for field access/unification
		table.DefineTypeAlias(stmt.Name.Value, tCon, realType, origin)
	} else {
		// ADT: Register constructors
		for _, c := range stmt.Constructors {
			// Check naming convention (constructor must start with uppercase)
			if !checkTypeName(c.Name.Value, c.Name.GetToken(), &errors) {
				continue
			}
			
			var resultType typesystem.Type = typesystem.TCon{Name: stmt.Name.Value}
			if len(stmt.TypeParameters) > 0 {
				args := []typesystem.Type{}
				for _, tp := range stmt.TypeParameters {
					args = append(args, typesystem.TVar{Name: tp.Value})
				}
				resultType = typesystem.TApp{Constructor: resultType, Args: args}
			}

			var constructorType typesystem.Type

			if len(c.Parameters) > 0 {
				var params []typesystem.Type
				for _, p := range c.Parameters {
					// Use typeScope to resolve type parameters
					params = append(params, BuildType(p, typeScope, &errors))
				}
				constructorType = typesystem.TFunc{
					Params:     params,
					ReturnType: resultType,
				}
			} else {
				constructorType = resultType
			}

			if table.IsDefined(c.Name.Value) {
				errors = append(errors, diagnostics.NewError(
					diagnostics.ErrA004,
					c.Name.GetToken(),
					c.Name.Value,
				))
				continue
			}
			table.DefineConstructor(c.Name.Value, constructorType, origin)
			table.RegisterVariant(stmt.Name.Value, c.Name.Value)

			for _, p := range c.Parameters {
				if nt, ok := p.(*ast.NamedType); ok {
					name := nt.Name.Value
					if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
						// Check if defined in typeScope (which includes global via outer)
						if !typeScope.IsDefined(name) {
							errors = append(errors, diagnostics.NewError(
								diagnostics.ErrA002,
								nt.GetToken(),
								name,
							))
						}
					}
				}
			}
		}
	}
	return errors
}

func (w *walker) VisitTypeDeclarationStatement(stmt *ast.TypeDeclarationStatement) {
	// Use RegisterTypeDeclaration to register and get errors
	errs := RegisterTypeDeclaration(stmt, w.symbolTable, w.currentModuleName)
	w.addErrors(errs)
}

func (w *walker) VisitTraitDeclaration(n *ast.TraitDeclaration) {
	// Skip if trait was not properly parsed
	if n == nil || n.Name == nil {
		return
	}
	
	// 0. Check naming convention (trait name must start with uppercase)
	if !checkTypeName(n.Name.Value, n.Token, &w.errors) {
		return
	}
	
	// 0.1. Check for redefinition of existing trait (including built-ins)
	if sym, ok := w.symbolTable.Find(n.Name.Value); ok && sym.Kind == symbols.TraitSymbol {
		w.addError(diagnostics.NewError(
			diagnostics.ErrA004,
			n.Token,
			n.Name.Value,
		))
		return
	}
	
	// 1. Extract super trait names and verify they exist
	superTraitNames := make([]string, 0, len(n.SuperTraits))
	for _, st := range n.SuperTraits {
		var superName string
		if nt, ok := st.(*ast.NamedType); ok {
			superName = nt.Name.Value
		}
		if superName != "" {
			// Check that super trait exists
			if !w.symbolTable.IsDefined(superName) {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA006,
					n.Token,
					superName,
				))
			} else {
				superTraitNames = append(superTraitNames, superName)
			}
		}
	}

	// 2. Register Trait with type params and super traits
	typeParamNames := make([]string, len(n.TypeParams))
	for i, tp := range n.TypeParams {
		typeParamNames[i] = tp.Value
	}
	w.symbolTable.DefineTrait(n.Name.Value, typeParamNames, superTraitNames, w.currentModuleName)

	// 3. Register methods
	// Methods are generic functions where the TypeParam is the trait variable.
	// e.g. show(val: a) -> String. 'a' is bound to the trait param.

	// We need a scope for the trait definition to define 'a'
	outer := w.symbolTable
	w.symbolTable = symbols.NewEnclosedSymbolTable(outer)
	defer func() { w.symbolTable = outer }()

	// Define the type variables
	for _, tp := range n.TypeParams {
		w.symbolTable.DefineType(tp.Value, typesystem.TVar{Name: tp.Value}, "")
	}

	for _, method := range n.Signatures {
		// Construct function type
		var retType typesystem.Type
		if method.ReturnType != nil {
			retType = BuildType(method.ReturnType, w.symbolTable, &w.errors)
		} else {
			retType = typesystem.Nil
		}

		var params []typesystem.Type
		for _, p := range method.Parameters {
			if p.Type != nil {
				params = append(params, BuildType(p.Type, w.symbolTable, &w.errors))
			} else {
				// Error: method signature usually requires types
				params = append(params, w.freshVar())
			}
		}

		methodType := typesystem.TFunc{
			Params:     params,
			ReturnType: retType,
		}

		// Register in Global Scope (outer) so it can be called
		// And associate with Trait
		outer.RegisterTraitMethod(method.Name.Value, n.Name.Value, methodType, w.currentModuleName)
		
		// Register method name in trait's method list
		outer.RegisterTraitMethod2(n.Name.Value, method.Name.Value)
		
		// If this is an operator method, register the operator -> trait mapping
		if method.Operator != "" {
			// Check if operator is already defined in another trait
			if existingTrait, exists := outer.GetTraitForOperator(method.Operator); exists {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA004,
					method.Token,
					"operator "+method.Operator+" (already defined in trait "+existingTrait+")",
				))
			} else {
				outer.RegisterOperatorTrait(method.Operator, n.Name.Value)
			}
		}
		
		// If method has a body, it's a default implementation
		if method.Body != nil {
			outer.RegisterTraitDefaultMethod(n.Name.Value, method.Name.Value)
			// Store the function for evaluator
			key := n.Name.Value + "." + method.Name.Value
			w.TraitDefaults[key] = method
		}
	}
}

func resolveReceiverTypeName(t ast.Type, table *symbols.SymbolTable) string {
	switch t := t.(type) {
	case *ast.NamedType:
		// Always return the named type - even if it resolves to a record (type alias)
		return t.Name.Value
	case *ast.TupleType:
		return "TUPLE"
	case *ast.RecordType:
		return "RECORD"
	default:
		return ""
	}
}

func (w *walker) VisitInstanceDeclaration(n *ast.InstanceDeclaration) {
	// Skip if not properly parsed
	if n == nil || n.TraitName == nil {
		return
	}
	
	// 1. Check if Trait exists
	if !w.symbolTable.IsDefined(n.TraitName.Value) {
		w.addError(diagnostics.NewError(
			diagnostics.ErrA001, // Undeclared identifier
			n.TraitName.GetToken(),
			n.TraitName.Value,
		))
		return
	}

	// 1b. Check that super traits are implemented for the target type
	// We need to build target type first to check implementations
	if n.Target == nil {
		w.addError(diagnostics.NewError(
			diagnostics.ErrA003,
			n.Token,
			"missing target type for instance",
		))
		return
	}
	targetType := BuildType(n.Target, w.symbolTable, &w.errors)
	
	// Kind check: For HKT traits like Functor<F>, F must be a type constructor
	// Use registered kinds (from symbol table) to verify
	// Automatically detect HKT traits by checking if type param is applied in method signatures
	isHKT := w.symbolTable.IsHKTTrait(n.TraitName.Value)
	
	// Extra type params in instance declaration count as partial application
	hasExtraParams := len(n.TypeParams) > 0
	
	if isHKT && !hasExtraParams {
		// Get type name directly from AST (before BuildType resolves aliases)
		var typeName string
		if namedType, ok := n.Target.(*ast.NamedType); ok {
			typeName = namedType.Name.Value
		}
		
		if typeName != "" {
			// Get target's kind from symbol table
			targetKind, hasKind := w.symbolTable.GetKind(typeName)
			
			isTypeConstructor := false
			if hasKind {
				// Kind * -> * or * -> * -> * etc means type constructor
				_, isArrow := targetKind.(typesystem.KArrow)
				isTypeConstructor = isArrow
			}
			
			if !isTypeConstructor {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA003,
					n.Target.GetToken(),
					fmt.Sprintf("type %s has kind *, but trait %s requires kind * -> * (type constructor)",
						typeName, n.TraitName.Value),
				))
				return
			}
		}
	}
	
	superTraits, _ := w.symbolTable.GetTraitSuperTraits(n.TraitName.Value)
	for _, superTrait := range superTraits {
		if !w.symbolTable.IsImplementationExists(superTrait, targetType) {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				n.Token,
				"cannot implement "+n.TraitName.Value+" for "+targetType.String()+": missing implementation of super trait "+superTrait,
			))
			return
		}
	}

	// 2. Extract type name from target
	var typeName string
	if tCon, ok := targetType.(typesystem.TCon); ok {
		typeName = tCon.Name
	} else if tApp, ok := targetType.(typesystem.TApp); ok {
		// Extract constructor name from app
		if tCon, ok := tApp.Constructor.(typesystem.TCon); ok {
			typeName = tCon.Name
		}
	}

	if typeName == "" {
		// Try to get from AST directly if BuildType resolves to something else (like Int which is built-in)
		if nt, ok := n.Target.(*ast.NamedType); ok {
			typeName = nt.Name.Value
		} else if _, ok := n.Target.(*ast.TupleType); ok {
			// Tuple support: use standardized "TUPLE" name
			typeName = "TUPLE"
		} else if _, ok := n.Target.(*ast.RecordType); ok {
			// Record support
			typeName = "RECORD"
		} else if _, ok := n.Target.(*ast.FunctionType); ok {
			// Function support
			typeName = "FUNCTION"
		}
	}

	if typeName == "" {
		// Fallback or error
		w.addError(diagnostics.NewError(
			diagnostics.ErrA003,
			n.Target.GetToken(),
			"invalid target type for instance",
		))
		return
	}

	// 3. Register Implementation
	err := w.symbolTable.RegisterImplementation(n.TraitName.Value, targetType)
	if err != nil {
		w.addError(diagnostics.NewError(
			diagnostics.ErrA004, // Redefinition/Overlap
			n.TraitName.GetToken(),
			err.Error(),
		))
		return
	}

	// 3b. Check that all required methods are implemented
	requiredMethods := w.symbolTable.GetTraitRequiredMethods(n.TraitName.Value)
	implementedMethods := make(map[string]bool)
	for _, method := range n.Methods {
		// For operator methods, Name is nil and Operator contains the operator symbol
		if method.Name != nil {
			implementedMethods[method.Name.Value] = true
		} else if method.Operator != "" {
			implementedMethods["("+method.Operator+")"] = true
		}
	}
	for _, required := range requiredMethods {
		if !implementedMethods[required] {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				n.Token,
				"instance "+n.TraitName.Value+" for "+typeName+" is missing required method '"+required+"'",
			))
		}
	}

	// 4. Analyze methods
	// Create a new scope for the implementation to avoid polluting global scope
	outer := w.symbolTable
	w.symbolTable = symbols.NewEnclosedSymbolTable(outer)
	defer func() { w.symbolTable = outer }()

	// Verify signatures
	typeParamNames, ok := w.symbolTable.GetTraitTypeParams(n.TraitName.Value)
	if !ok {
		w.addError(diagnostics.NewError(
			diagnostics.ErrA001,
			n.TraitName.GetToken(),
			"unknown trait type param for "+n.TraitName.Value,
		))
		return
	}

	// Create substitution: TraitTypeParam[0] -> TargetType
	// Important: If targetType contains type variables with the same name as the trait's
	// type parameter (e.g., instance UserOpChoose Box<T> where trait has parameter T),
	// we need to rename them to avoid infinite substitution.
	// The targetType's free type variables are instance type parameters, not the trait's.
	renamedTarget := renameConflictingTypeVars(targetType, typeParamNames, w.inferCtx)
	subst := typesystem.Subst{typeParamNames[0]: renamedTarget}

	for _, method := range n.Methods {
		method.Accept(w)

		// Verify signature matches Trait definition
		// 1. Find generic signature
		genericSymbol, ok := w.symbolTable.Find(method.Name.Value)
		if !ok {
			// Method not in trait?
			w.addError(diagnostics.NewError(
				diagnostics.ErrA001,
				method.Name.GetToken(),
				"method "+method.Name.Value+" is not part of trait "+n.TraitName.Value,
			))
			continue
		}

		traitForMethod, _ := w.symbolTable.GetTraitForMethod(method.Name.Value)
		if traitForMethod != n.TraitName.Value {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				method.Name.GetToken(),
				"method "+method.Name.Value+" belongs to trait "+traitForMethod+", not "+n.TraitName.Value,
			))
			continue
		}

		// 2. Create Expected Type (Generic signature with substitution)
		expectedType := genericSymbol.Type.Apply(subst)

		// 3. Get Actual Type from the method definition in CURRENT scope
		// VisitFunctionStatement (called by method.Accept(w)) defined it in w.symbolTable (the inner scope).
		actualSymbol, ok := w.symbolTable.Find(method.Name.Value)
		if !ok {
			// Should not happen if VisitFunctionStatement works
			continue
		}
		actualType := actualSymbol.Type

		// 4. Unify
		_, err := typesystem.Unify(expectedType, actualType)
		if err != nil {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				method.Name.GetToken(),
				"method signature mismatch: expected "+expectedType.String()+", got "+actualType.String(),
			))
		} else {
			// 5. Register instance method signature for use in type inference
			// This allows traits like Optional to correctly extract inner types
			// for any user-defined type, not just built-in types.
			outer.RegisterInstanceMethod(n.TraitName.Value, typeName, method.Name.Value, expectedType)
		}
	}
}

func (w *walker) VisitConstantDeclaration(n *ast.ConstantDeclaration) {
	// Skip if not properly parsed
	if n == nil {
		return
	}
	
	// Handle pattern destructuring: (a, b) :- pair
	if n.Pattern != nil {
		w.visitPatternDeclaration(n)
		return
	}
	
	// Simple binding: x :- expr
	if n.Name == nil {
		return
	}
	
	// 0. Check naming convention (must start with lowercase)
	if !checkValueName(n.Name.Value, n.Name.GetToken(), &w.errors) {
		return
	}
	
	// 1. Check for redefinition - constants can NEVER be redefined
	if w.symbolTable.IsDefined(n.Name.Value) {
		sym, ok := w.symbolTable.Find(n.Name.Value)
		// Only allow if it's a pending symbol (forward declaration)
		if ok && !sym.IsPending {
			w.addError(diagnostics.NewError(diagnostics.ErrA004, n.Name.GetToken(), n.Name.Value))
			return
		}
	}

	// 2. Infer Value Type
	valType, s1, err := InferWithContext(w.inferCtx, n.Value, w.symbolTable)
	if err != nil {
		w.appendError(n.Value, err)
		valType = w.freshVar() // Recovery
	}
	valType = valType.Apply(s1)

	// 3. Check Type Annotation (if present)
	if n.TypeAnnotation != nil {
		annotType := BuildType(n.TypeAnnotation, w.symbolTable, &w.errors)
		subst, err := typesystem.Unify(annotType, valType)
		if err != nil {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				n.Value.GetToken(),
				"constant type mismatch: expected "+annotType.String()+", got "+valType.String(),
			))
		}
		// Use the annotated type as the source of truth
		valType = annotType
		valType = valType.Apply(subst)
	}

	// 4. Register in Symbol Table as Constant (immutable)
	w.symbolTable.DefineConstant(n.Name.Value, valType, w.currentModuleName)
}

// visitPatternDeclaration handles pattern destructuring in constant bindings (:-) 
func (w *walker) visitPatternDeclaration(n *ast.ConstantDeclaration) {
	// 1. Infer Value Type
	valType, s1, err := InferWithContext(w.inferCtx, n.Value, w.symbolTable)
	if err != nil {
		w.appendError(n.Value, err)
		return
	}
	valType = valType.Apply(s1)
	
	// 2. Bind pattern variables with inferred types as CONSTANTS
	w.bindPatternVariablesAsConstant(n.Pattern, valType, n.Token)
}

// bindPatternVariables binds variables from a pattern to their types (mutable)
func (w *walker) bindPatternVariables(pat ast.Pattern, valType typesystem.Type, tok token.Token) {
	w.bindPatternVariablesWithConstFlag(pat, valType, tok, false)
}

// bindPatternVariablesAsConstant binds variables from a pattern to their types (immutable)
func (w *walker) bindPatternVariablesAsConstant(pat ast.Pattern, valType typesystem.Type, tok token.Token) {
	w.bindPatternVariablesWithConstFlag(pat, valType, tok, true)
}

// bindPatternVariablesWithConstFlag binds variables from a pattern to their types
func (w *walker) bindPatternVariablesWithConstFlag(pat ast.Pattern, valType typesystem.Type, tok token.Token, isConstant bool) {
	switch p := pat.(type) {
	case *ast.IdentifierPattern:
		// Check naming convention (variable must start with lowercase)
		if !checkValueName(p.Value, p.Token, &w.errors) {
			return
		}
		// Check for redefinition
		if w.symbolTable.IsDefined(p.Value) {
			sym, ok := w.symbolTable.Find(p.Value)
			if ok && !sym.IsPending {
				w.addError(diagnostics.NewError(diagnostics.ErrA004, p.Token, p.Value))
				return
			}
		}
		if isConstant {
			w.symbolTable.DefineConstant(p.Value, valType, "")
		} else {
			w.symbolTable.Define(p.Value, valType, "")
		}
		
	case *ast.TuplePattern:
		// valType should be TTuple
		if tuple, ok := valType.(typesystem.TTuple); ok {
			if len(tuple.Elements) != len(p.Elements) {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA003,
					tok,
					fmt.Sprintf("tuple pattern has %d elements but value has %d", len(p.Elements), len(tuple.Elements)),
				))
				return
			}
			for i, elem := range p.Elements {
				w.bindPatternVariablesWithConstFlag(elem, tuple.Elements[i], tok, isConstant)
			}
		} else {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				tok,
				"cannot destructure non-tuple value with tuple pattern",
			))
		}
		
	case *ast.ListPattern:
		// valType should be TApp List<T>
		if app, ok := valType.(typesystem.TApp); ok {
			if con, ok := app.Constructor.(typesystem.TCon); ok && con.Name == "List" && len(app.Args) > 0 {
				elemType := app.Args[0]
				for _, elem := range p.Elements {
					w.bindPatternVariablesWithConstFlag(elem, elemType, tok, isConstant)
				}
			} else {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA003,
					tok,
					"cannot destructure non-list value with list pattern",
				))
			}
		} else {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				tok,
				"cannot destructure non-list value with list pattern",
			))
		}
		
	case *ast.WildcardPattern:
		// Ignore - don't bind anything
		
	case *ast.RecordPattern:
		// Handle both TRecord and named record types
		var fields map[string]typesystem.Type
		
		switch t := valType.(type) {
		case typesystem.TRecord:
			fields = t.Fields
		default:
			// Try to get underlying type if it's a named record type
			if underlying := typesystem.UnwrapUnderlying(valType); underlying != nil {
				if rec, ok := underlying.(typesystem.TRecord); ok {
					fields = rec.Fields
				}
			}
		}
		
		if fields == nil {
			w.addError(diagnostics.NewError(
				diagnostics.ErrA003,
				tok,
				"cannot destructure non-record value with record pattern",
			))
			return
		}
		
		for fieldName, fieldPat := range p.Fields {
			fieldType, ok := fields[fieldName]
			if !ok {
				w.addError(diagnostics.NewError(
					diagnostics.ErrA003,
					tok,
					fmt.Sprintf("record does not have field '%s'", fieldName),
				))
				return
			}
			w.bindPatternVariablesWithConstFlag(fieldPat, fieldType, tok, isConstant)
		}
		
	default:
		w.addError(diagnostics.NewError(
			diagnostics.ErrA003,
			tok,
			"unsupported pattern in destructuring",
		))
	}
}

func (w *walker) VisitNamedType(n *ast.NamedType) {}

func (w *walker) VisitDataConstructor(n *ast.DataConstructor) {}

func (w *walker) VisitTupleType(t *ast.TupleType) {}

func (w *walker) VisitFunctionType(n *ast.FunctionType) {
	// Just check sub types
	for _, p := range n.Parameters {
		p.Accept(w)
	}
	n.ReturnType.Accept(w)
}

func (w *walker) VisitRecordType(n *ast.RecordType) {
	for _, v := range n.Fields {
		v.Accept(w)
	}
}

func (w *walker) VisitUnionType(n *ast.UnionType) {
	for _, t := range n.Types {
		t.Accept(w)
	}
}

// renameConflictingTypeVars renames type variables in `t` that conflict with `conflictNames`.
// This is needed when creating substitutions for trait instances where the target type
// might have type variables with the same name as the trait's type parameters.
// For example: `instance UserOpChoose Box<T>` where trait UserOpChoose<T> - the T in Box<T>
// should not be confused with the trait's T parameter.
func renameConflictingTypeVars(t typesystem.Type, conflictNames []string, ctx *InferenceContext) typesystem.Type {
	if t == nil || ctx == nil {
		return t
	}
	
	// Find conflicting type variables in t
	freeVars := t.FreeTypeVariables()
	conflictSet := make(map[string]bool)
	for _, name := range conflictNames {
		conflictSet[name] = true
	}
	
	// Create renaming substitution for conflicting vars
	renameSubst := typesystem.Subst{}
	for _, tv := range freeVars {
		if conflictSet[tv.Name] {
			// Rename to a fresh variable
			fresh := ctx.FreshVar()
			renameSubst[tv.Name] = fresh
		}
	}
	
	if len(renameSubst) == 0 {
		return t // No conflicts
	}
	
	return t.Apply(renameSubst)
}

// importTraitImplementations copies trait implementations for exported types
// from imported module's symbol table to current symbol table
func (w *walker) importTraitImplementations(importedTable *symbols.SymbolTable, exportedTypes map[string]bool, moduleName string) {
	// Get all implementations from imported module
	allImpls := importedTable.GetAllImplementations()

	// For each trait, copy implementations for exported types
	for traitName, impls := range allImpls {
		for _, implType := range impls {
			// Check if this implementation should be imported
			// We import implementations for types that are either:
			// 1. Exported types (by name)
			// 2. Types that match exported type structures (for aliases)

			if w.shouldImportTraitImplementation(implType, exportedTypes, moduleName) {
				// Tag the type with module name if it's a named type
				taggedType := tagModule(implType, moduleName, exportedTypes)
				// Register the implementation in current table
				_ = w.symbolTable.RegisterImplementation(traitName, taggedType)
			}
		}
	}
}

// shouldImportTraitImplementation determines if a trait implementation should be imported
func (w *walker) shouldImportTraitImplementation(implType typesystem.Type, exportedTypes map[string]bool, moduleName string) bool {
	switch t := implType.(type) {
	case typesystem.TCon:
		// Import if the type name is exported
		return exportedTypes[t.Name]
	case typesystem.TRecord:
		// Import record types (for type aliases that resolve to records)
		return true
	case typesystem.TApp:
		// For type applications (like List<T>), check the constructor
		if tCon, ok := t.Constructor.(typesystem.TCon); ok {
			return exportedTypes[tCon.Name]
		}
	}

	// Default: don't import
	return false
}

// importExtensionMethods copies extension methods for exported types
// from imported module's symbol table to current symbol table
func (w *walker) importExtensionMethods(importedTable *symbols.SymbolTable, exportedTypes map[string]bool) {
	allExtMethods := importedTable.GetAllExtensionMethods()
	
	for typeName, methods := range allExtMethods {
		// Only import extension methods for exported types
		if !exportedTypes[typeName] {
			continue
		}
		
		for methodName, methodType := range methods {
			w.symbolTable.RegisterExtensionMethod(typeName, methodName, methodType)
		}
	}
}
