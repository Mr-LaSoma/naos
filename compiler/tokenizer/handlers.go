package tokenizer

import "naoslang/tokenizer/tokens"

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

	// operations
	return t.newToken(tokens.TOKNotImplemented, true)
}
