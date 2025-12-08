package lexer

import (
	"github.com/funvibe/funxy/internal/pipeline"
	"github.com/funvibe/funxy/internal/token"
)

const lookaheadBufferSize = 10

type bufferedLexer struct {
	l      *Lexer
	buffer []token.Token
	pos    int
}

func NewTokenStream(l *Lexer) pipeline.TokenStream {
	return &bufferedLexer{l: l}
}

func newBufferedLexer(l *Lexer) *bufferedLexer {
	return &bufferedLexer{l: l}
}

func (bl *bufferedLexer) Next() token.Token {
	if bl.pos < len(bl.buffer) {
		tok := bl.buffer[bl.pos]
		bl.pos++
		return tok
	}

	tok := bl.l.NextToken()
	return tok
}

func (bl *bufferedLexer) Peek(n int) []token.Token {
	// Ensure the buffer has at least one token if it's currently empty,
	// to prevent panics on Peek(0) at the end of the stream.
	if len(bl.buffer)-bl.pos == 0 {
		nextTok := bl.l.NextToken()
		bl.buffer = append(bl.buffer, nextTok)
		if nextTok.Type == token.EOF {
			// If we hit EOF, we don't need to read more.
			// This slice will be returned below.
		}
	}

	// Ensure buffer has enough tokens for the requested lookahead
	for len(bl.buffer)-bl.pos < n {
		nextTok := bl.l.NextToken()
		bl.buffer = append(bl.buffer, nextTok)
		if nextTok.Type == token.EOF {
			break
		}
	}

	// Trim buffer if it's too large
	if bl.pos > lookaheadBufferSize {
		bl.buffer = bl.buffer[bl.pos:]
		bl.pos = 0
	}

	end := bl.pos + n
	if end > len(bl.buffer) {
		end = len(bl.buffer)
	}

	return bl.buffer[bl.pos:end]
}

var _ pipeline.TokenStream = (*bufferedLexer)(nil)

type LexerProcessor struct{}

func (lp *LexerProcessor) Process(ctx *pipeline.PipelineContext) *pipeline.PipelineContext {
	l := New(ctx.SourceCode)
	stream := newBufferedLexer(l)
	ctx.TokenStream = stream
	return ctx
}
