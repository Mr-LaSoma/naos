package tokens

import "fmt"

// +-------------+
// | Parentheses |
// +-------------+

var parentheses = map[rune]TokenKind{
	'(': TOKLParen, ')': TOKRParen, '[': TOKLBrackets,
	']': TOKRBrackets, '{': TOKLBraces, '}': TOKRBraces,
}

func RuneIsParentheses(ch rune) bool {
	_, ok := parentheses[ch]
	return ok
}

func parenthesisToKind(ch rune) (TokenKind, error) {
	kind, ok := parentheses[ch]
	if !ok {
		return 0, fmt.Errorf("%c is not a valid parentheses", ch)
	}

	return kind, nil
}
