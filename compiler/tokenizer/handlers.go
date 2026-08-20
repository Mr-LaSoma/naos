package tokenizer

import (
	"fmt"
	"naoslang/tokenizer/tokens"
	"unicode"
)

// +----------+
// | Comments |
// +----------+

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

// +-----------+
// | Operators |
// +-----------+

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

// +------------------------+
// | Identifiers / Keywords |
// +------------------------+

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

// +---------+
// | Numbers |
// +---------+

func (t *tokenizer) handleNumbers() tokens.Token {
	ch := t.peekRune()
	base := tokens.NumericBaseDecimal

	if t.getTokenValue() == "0" {
		nbase, err := tokens.RuneToNumericBase(ch)
		if err == nil {
			base = nbase
			tokens.PositionGoForward(&t.currentPos, false)
			ch = t.peekRune()
		}
	}

	for tokens.IsHexDigit(ch) || ch == '_' {
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	if ch == '.' {
		tokens.PositionGoForward(&t.currentPos, false)
		return t.handleFloats()
	}

	if base != tokens.NumericBaseDecimal && len(t.getTokenValue()) < 3 { // the bases must be at least 3 long (0xF -> valid, 0x -> invalid)
		return t.newInvalidToken(fmt.Sprintf("invalid number literal, %v base digit must be at least 3 characters long", base))
	}

	str := t.getTokenValue()
	digits := t.getTokenValue()
	wasUnderscore := false
	if base != tokens.NumericBaseDecimal {
		digits = digits[2:]
	}

	if len(digits) > 0 && digits[0] == '_' {
		return t.newInvalidToken(fmt.Sprintf("invalid number literal %s, digits cannot start with an underscore", str))
	}

	for _, digit := range digits {
		if digit == '_' && wasUnderscore {
			return t.newInvalidToken(fmt.Sprintf("invalid number literal %s, can't have two _ sequentially", str))
		}

		if !tokens.IsValidBaseDigit(digit, base) && digit != '_' {
			return t.newInvalidToken(fmt.Sprintf("invalid number literal %s, unexpected %v base digit found: %c", str, base, digit))
		}
		wasUnderscore = digit == '_'
	}

	if str[len(str)-1] == '_' { // last character of te digit was a _
		return t.newInvalidToken(fmt.Sprintf("invalid number literal %s, digit can't have an _ as last character", str))
	}

	return t.newToken(tokens.TOKInt, true)
}

func (t *tokenizer) handleFloats() tokens.Token {
	ch := t.peekRune()

	for tokens.IsDecimalDigit(ch) || ch == '_' || ch == '.' {
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	digits := t.getTokenValue()
	gotPeriod := false
	wasPeriod := false
	wasUnderscore := false
	for _, digit := range digits {
		if digit == '.' {
			if gotPeriod {
				return t.newInvalidToken(fmt.Sprintf("invalid float literal %s, cannot have multiple decimal points", digits))
			}
			gotPeriod = true
		}

		if digit == '_' {
			if wasUnderscore {
				return t.newInvalidToken(fmt.Sprintf("invalid float literal %s, can't have two _ sequentially", digits))
			}
			if wasPeriod {
				return t.newInvalidToken(fmt.Sprintf("invalid float literal %s, can't have _ after .", digits))
			}
		}

		if !tokens.IsDecimalDigit(digit) && digit != '_' && digit != '.' {
			return t.newInvalidToken(fmt.Sprintf("invalid float literal %s, expected decimal digit found %c", digits, digit))
		}

		wasPeriod = digit == '.'
		wasUnderscore = digit == '_'
	}

	if digits[len(digits)-1] == '_' { // last character of te digit was a _
		return t.newInvalidToken(fmt.Sprintf("invalid float literal %s, digit can't have an _ as last character", digits))
	}

	return t.newToken(tokens.TOKFloat, true)
}

// +---------+
// | Strings |
// +---------+

func (t *tokenizer) handleEscapeSequences(isChar bool) (bool, string) {
	ch := t.peekRune()
	if tokens.IsBaseEscapeSeq(ch) {
		tokens.PositionGoForward(&t.currentPos, false)
		if isChar && ch == '"' {
			return false, `\"`
		}
		if !isChar && ch == '\'' {
			return false, `\'`
		}

		return true, `\` + string(ch)
	}

	seq := []rune{'\\'}
	kind, err := tokens.RuneToEscapeKind(ch)
	if err != nil {
		if tokens.IsOctalDigit(ch) {
			kind = tokens.EscapeKindOctal
		} else {
			return false, `\` + string(ch)
		}
	} else {
		seq = append(seq, ch)
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	digitsNeeded := kind.NumbersNeeded()
	base := kind.ToNumericBase()
	for digitsNeeded > 0 && tokens.IsValidBaseDigit(ch, base) {
		seq = append(seq, ch)

		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
		digitsNeeded--
	}

	if digitsNeeded > 0 && kind != tokens.EscapeKindOctal {
		return false, string(seq)
	}

	return true, string(seq)
}

func (t *tokenizer) handleStrings() tokens.Token {
	ch := t.peekRune()
	wasValidEscape := true
	escapeSeq := ""
	gotNewLine := false
	for ch != tokens.EOFRune {
		if ch == '\\' {
			tokens.PositionGoForward(&t.currentPos, false)
			if !wasValidEscape {
				_, _ = t.handleEscapeSequences(false)
			} else {
				wasValidEscape, escapeSeq = t.handleEscapeSequences(false)
			}

			ch = t.peekRune()
			continue
		}

		if ch == '"' {
			tokens.PositionGoForward(&t.currentPos, false)
			if gotNewLine {
				return t.newInvalidToken(fmt.Sprintf("invalid string literal %s, string was never closed", t.getTokenValue()))
			}
			if !wasValidEscape {
				return t.newInvalidToken(fmt.Sprintf("invalid string literal %s, %s is not a valid escape sequence", t.getTokenValue(), escapeSeq))
			}
			return t.newToken(tokens.TOKString, true)
		}

		if ch == '\n' {
			gotNewLine = true
			tokens.PositionGoForward(&t.currentPos, true)
			ch = t.peekRune()
			continue
		}
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	return t.newInvalidToken(fmt.Sprintf("invalid string literal %s, expected \" found EOF", t.getTokenValue()))
}

func (t *tokenizer) handleChars() tokens.Token {
	ch := t.peekRune()
	wasValidEscape := true
	escapeSeq := ""
	nChars := 0
	gotNewLine := false

	for ch != tokens.EOFRune {
		if ch == '\\' {
			tokens.PositionGoForward(&t.currentPos, false)
			nChars++

			if !wasValidEscape {
				_, _ = t.handleEscapeSequences(true)
			} else {
				wasValidEscape, escapeSeq = t.handleEscapeSequences(true)
			}

			ch = t.peekRune()
			continue
		}

		if ch == '\'' {
			tokens.PositionGoForward(&t.currentPos, false)
			if gotNewLine {
				return t.newInvalidToken(fmt.Sprintf("invalid char literal %s, char was never closed", t.getTokenValue()))
			}
			if nChars > 1 {
				return t.newInvalidToken(fmt.Sprintf("invalid char literal %s, chars can't have more then 1 char", t.getTokenValue()))
			}
			if nChars == 0 {
				return t.newInvalidToken(fmt.Sprintf("invalid char literal %s, chars can't be empty", t.getTokenValue()))
			}
			if !wasValidEscape {
				return t.newInvalidToken(fmt.Sprintf("invalid char literal %s, %s is not a valid escape sequence", t.getTokenValue(), escapeSeq))
			}
			return t.newToken(tokens.TOKChar, true)
		}

		if ch == '\n' {
			gotNewLine = true
			tokens.PositionGoForward(&t.currentPos, true)
			ch = t.peekRune()
			continue
		}

		nChars++
		tokens.PositionGoForward(&t.currentPos, false)
		ch = t.peekRune()
	}

	return t.newInvalidToken(fmt.Sprintf("invalid char literal %s, expected \" found EOF", t.getTokenValue()))
}
