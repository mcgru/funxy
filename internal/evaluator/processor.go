package evaluator

import (
	"github.com/funvibe/funxy/internal/diagnostics"
	"github.com/funvibe/funxy/internal/modules"
	"github.com/funvibe/funxy/internal/pipeline"
	"github.com/funvibe/funxy/internal/token"
	"path/filepath"
)

type EvaluatorProcessor struct{}

func (ep *EvaluatorProcessor) Process(ctx *pipeline.PipelineContext) *pipeline.PipelineContext {
	if ctx.AstRoot == nil || len(ctx.Errors) > 0 {
		return ctx
	}

	eval := New()
	// Use shared loader from analyzer (with type assertion)
	if loader, ok := ctx.Loader.(*modules.Loader); ok {
		eval.SetLoader(loader)
	} else {
		eval.SetLoader(modules.NewLoader())
	}
	eval.TraitDefaults = ctx.TraitDefaults   // Pass trait defaults from analyzer
	eval.OperatorTraits = ctx.OperatorTraits // Pass operator -> trait mappings
	eval.TypeMap = ctx.TypeMap               // Pass inferred types from analyzer
	
	// Set BaseDir and CurrentFile from ctx.FilePath
	if ctx.FilePath != "" {
		dir := filepath.Dir(ctx.FilePath)
		eval.BaseDir = dir
		eval.CurrentFile = filepath.Base(ctx.FilePath)
	} else {
		eval.BaseDir = "."
		eval.CurrentFile = "<stdin>"
	}

	env := NewEnvironment()
	RegisterBuiltins(env)
	RegisterFPTraits(eval, env) // Register FP traits (Semigroup, Monoid, Functor, Applicative, Monad)
	eval.GlobalEnv = env // Store for default implementations

	result := eval.Eval(ctx.AstRoot, env)
	if result != nil && result.Type() == ERROR_OBJ {
		// Convert evaluator error to diagnostic error
		// We don't have token info in result directly easily, unless we pass it.
		// For now, use a dummy token or modify Evaluator to return Error with Token?
		// Evaluator Errors are runtime errors, usually not attached to source code location unless we track it.
		ctx.Errors = append(ctx.Errors, diagnostics.NewError(
			diagnostics.ErrR001,
			token.Token{}, // Missing token info
			result.Inspect(),
		))
	}

	return ctx
}

