package tokenizer

import (
	"fmt"
	"naoslang/tokenizer/tokens"
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
	operator := []rune{start}
	ch := t.peekRune()
	for tokens.StringIsOperator(string(append(operator, ch))) {
		operator = append(operator, ch)
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	kind, err := tokens.OperatorToKind(string(operator))
	if err != nil {
		panic(fmt.Sprintf("unexpected error while lexing an operator: %v", err))
	}
	return t.newToken(kind, false)
}
