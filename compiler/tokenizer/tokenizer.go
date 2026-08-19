package tokenizer

import (
	"naoslang/tokenizer/tokens"
)

type tokenizer struct {
	source     []rune
	currentPos tokens.Position

	startTokPos tokens.Position
}

func New(source, filename string) *tokenizer {
	return &tokenizer{
		source:     []rune(source),
		currentPos: tokens.NewPosition(filename),
	}
}

func (t *tokenizer) NextToken() tokens.Token {
	t.skipWhiteSpace()

	t.startTokPos = t.currentPos
	switch t.readRune() {
	case '/':
		return t.handleComments()

	case tokens.EOFRune:
		return t.newToken(tokens.TOKEof, false)
	}

	return t.newToken(tokens.TOKNotImplemented, true)
}

func (t *tokenizer) Tokenize() []tokens.Token {
	var toks []tokens.Token
	var last tokens.Token

	for last.Kind != tokens.TOKEof {
		last = t.NextToken()
		toks = append(toks, last)
	}
	return toks
}
