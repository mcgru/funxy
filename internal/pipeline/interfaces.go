package pipeline

import (
	"github.com/funvibe/funxy/internal/token"
)

// Processor is any component that can process a
// PipelineContext and return a modified context.
type Processor interface {
	Process(ctx *PipelineContext) *PipelineContext
}

// TokenStream defines the contract for a buffered token stream.
type TokenStream interface {
	// Next consumes and returns the next token from the stream.
	Next() token.Token

	// Peek returns the next n tokens without consuming them.
	// If the stream has fewer than n tokens, it returns all remaining tokens.
	Peek(n int) []token.Token
}
