package tokenizer

import (
	"fmt"
	"naoslang/tokenizer/tokens"
	"unicode"
)

// handleComments handles comments tokens and it dispatches
// the / to the operations if is not a comment
func (t *tokenizer) handleComments() tokens.Token {
	switch t.peekRune() {
	case '/':
		tokens.PositionGoForward(&t.currentPos, false)
		t.skipMonoComment()
		return t.newToken(tokens.TOKMonoComment, false)
	case '*':
		tokens.PositionGoForward(&t.currentPos, false)
		t.skipMultiComment()
		return t.newToken(tokens.TOKMultiComment, false)
	}

	return t.handleOperators('/')
}

func (t *tokenizer) handleOperators(start rune) tokens.Token {
	ch := t.peekRune()
	for tokens.StringIsOperator(t.getTokenValue() + string(ch)) {
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	kind, err := tokens.OperatorToKind(t.getTokenValue())
	if err != nil {
		panic(fmt.Sprintf("unexpected error while lexing an operator: %v", err))
	}
	return t.newToken(kind, false)
}

func (t *tokenizer) handleIdentifiers() tokens.Token {
	ch := t.peekRune()
	for unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	if tokens.StringIsKeyword(t.getTokenValue()) {
		kind, err := tokens.KeywordToKind(t.getTokenValue())
		if err != nil {
			panic(fmt.Sprintf("unexpected error while lexing a keyword: %v", err))
		}
		return t.newToken(kind, false)
	}

	return t.newToken(tokens.TOKId, true)
}
