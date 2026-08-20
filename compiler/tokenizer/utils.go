package tokenizer

import (
	"naoslang/tokenizer/tokens"
	"unicode"
)

// newToken returns a new token created with a specific type.
// withLexeme signals if should be stored what is inside the token (true for literals, or identifiers)
func (t *tokenizer) newToken(kind tokens.TokenKind, withLexeme bool) tokens.Token {
	var s string

	if withLexeme { // used so i won't get the string for stuff like if, while ecc...
		s = t.getTokenValue()
	}
	return tokens.Token{
		Kind:   kind,
		Lexeme: s,
		Start:  t.startTokPos,
		End:    t.currentPos,
	}
}

func (t *tokenizer) newInvalidToken(errormsg string) tokens.Token {
	return tokens.Token{
		Kind:   tokens.TOKInvalid,
		Lexeme: errormsg,
		Start:  t.startTokPos,
		End:    t.currentPos,
	}
}

func (t *tokenizer) getTokenValue() string {
	return string(t.source[t.startTokPos.Offset:t.currentPos.Offset])
}

// readRune returns the current rune and goes forward with the position.
func (t *tokenizer) readRune() rune {
	if t.currentPos.Offset >= len(t.source) {
		return tokens.EOFRune
	}

	ch := t.source[t.currentPos.Offset]
	tokens.PositionGoForward(&t.currentPos, ch == '\n')
	return ch
}

// peekRune returns the current rune without moving.
func (t *tokenizer) peekRune() rune {
	if t.currentPos.Offset >= len(t.source) {
		return tokens.EOFRune
	}
	return t.source[t.currentPos.Offset]
}

// skipWhiteSpace skips all the white spaces
func (t *tokenizer) skipWhiteSpace() {
	for t.currentPos.Offset < len(t.source) {
		ch := t.source[t.currentPos.Offset]
		if unicode.IsSpace(ch) {
			tokens.PositionGoForward(&t.currentPos, ch == '\n')
			continue
		}
		break
	}
}

func (t *tokenizer) skipMonoComment() {
	for t.currentPos.Offset < len(t.source) {
		ch := t.source[t.currentPos.Offset]
		if ch == '\n' {
			break
		}
		tokens.PositionGoForward(&t.currentPos, ch == '\n')
	}
}

func (t *tokenizer) skipMultiComment() {
	nNested := 1
	ch := t.peekRune()
	for ch != tokens.EOFRune {
		tokens.PositionGoForward(&t.currentPos, ch == '\n')

		if ch == '/' && t.peekRune() == '*' {
			tokens.PositionGoForward(&t.currentPos, false) // skips the *
			nNested++
		}

		if ch == '*' && t.peekRune() == '/' {
			tokens.PositionGoForward(&t.currentPos, false) // skips the /
			nNested--
		}

		if nNested == 0 {
			break
		}
		ch = t.peekRune()
	}
}
