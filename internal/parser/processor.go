package parser

import (
	"github.com/funvibe/funxy/internal/diagnostics"
	"github.com/funvibe/funxy/internal/pipeline"
	"github.com/funvibe/funxy/internal/token"
)

type ParserProcessor struct{}

func (pp *ParserProcessor) Process(ctx *pipeline.PipelineContext) *pipeline.PipelineContext {
	if ctx.TokenStream == nil {
		// This case should ideally not be hit if lexer runs first, but as a safeguard:
		err := diagnostics.NewError("P000", token.Token{}, "parser: token stream is nil")
		ctx.Errors = append(ctx.Errors, err)
		return ctx
	}

	parser := New(ctx.TokenStream, ctx)
	ctx.AstRoot = parser.ParseProgram()

	// Errors are already added to the context by the parser instance,
	// so we don't need to retrieve them again.

	return ctx
}
