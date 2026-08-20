package tokenizer

import (
	"fmt"
	"naoslang/tokenizer/tokens"
	"unicode"
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
	ch := t.readRune()

	switch ch {
	case '\'':
		return t.newToken(tokens.TOKNotImplemented, true)
		panic("chars are not yet implemented")
	case '"':
		return t.newToken(tokens.TOKNotImplemented, true)
		panic("strings are not yet implemented")
	case '/':
		return t.handleComments()
	case tokens.EOFRune:
		return t.newToken(tokens.TOKEof, false)
	}

	if unicode.IsDigit(ch) {
		return t.newToken(tokens.TOKNotImplemented, true)
		panic("numbers are not yet implemented")
	}
	if unicode.IsLetter(ch) || ch == '_' {
		return t.handleIdentifiers()
	}

	if tokens.RuneIsParentheses(ch) {
		kind, err := tokens.ParenthesisToKind(ch)
		if err != nil {
			panic(fmt.Sprintf("unexpected error while lexing a parenthesis: %v", err))
		}
		return t.newToken(kind, false)
	}

	if tokens.RuneIsPunctuation(ch) {
		kind, err := tokens.PunctuationToKind(ch)
		if err != nil {
			panic(fmt.Sprintf("unexpected error while lexing a parenthesis: %v", err))
		}
		return t.newToken(kind, false)
	}

	if tokens.RuneIsOperatorStart(ch) {
		return t.handleOperators(ch)
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
