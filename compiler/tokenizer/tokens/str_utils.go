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

func ParenthesisToKind(ch rune) (TokenKind, error) {
	kind, ok := parentheses[ch]
	if !ok {
		return 0, fmt.Errorf("%c is not a valid parentheses", ch)
	}
	return kind, nil
}

// +-------------+
// | Punctuation |
// +-------------+

var punctuation = map[rune]TokenKind{';': TOKSemiColon, ',': TOKComma, '.': TOKPeriod}

func RuneIsPunctuation(ch rune) bool {
	_, ok := punctuation[ch]
	return ok
}

func PunctuationToKind(ch rune) (TokenKind, error) {
	kind, ok := punctuation[ch]
	if !ok {
		return 0, fmt.Errorf("%c is not a valid punctuation", ch)
	}
	return kind, nil
}

var (
	operators = map[string]TokenKind{
		"+": TOKPlus, "-": TOKMinus, "*": TOKStar, "/": TOKSlash,
		"%": TOKPercent, "=": TOKAssign, "+=": TOKPlusAssign,
		"-=": TOKMinusAssign, "*=": TOKStarAssign, "/=": TOKSlashAssign,
		"%=": TOKPercentAssign, "++": TOKIncrement, "--": TOKDecrement,

		"!": TOKNot, "&&": TOKAnd, "||": TOKOr,
		"<": TOKLess, ">": TOKGreater, "==": TOKEqual,
		"!=": TOKNotEqual, "<=": TOKLessEqual, ">=": TOKGreaterEqual,

		"&": TOKBwAnd, "|": TOKBwOr, "^": TOKBwXor,
		"~": TOKBwNot, "<<": TOKBwLShift, ">>": TOKBwRShift,
		"&=": TOKBwAndAssign, "|=": TOKBwOrAssign, "^=": TOKBwXorAssign,
		"<<=": TOKBwLShiftAssign, ">>=": TOKBwRShiftAssign,

		"#": TOKHashtag, "@": TOKAt, "->": TOKArrow,
		"=>": TOKBigArrow, ":": TOKColon, "::": TOKDoubleColon,
	}
	operatorsStart = "+-*/%=!&|<>^~#@:"
)

func RuneIsOperatorStart(ch rune) bool {
	for _, opStart := range operatorsStart {
		if ch == opStart {
			return true
		}
	}
	return false
}

func StringIsOperator(str string) bool {
	_, ok := operators[str]
	return ok
}

func OperatorToKind(str string) (TokenKind, error) {
	kind, ok := operators[str]
	if !ok {
		return 0, fmt.Errorf("%s is not a valid operator", str)
	}
	return kind, nil
}

// +----------+
// | Keywords |
// +----------+

var keywords = map[string]TokenKind{
	"package": TOKPackage, "using": TOKUsing,
	"pub": TOKPub, "priv": TOKPriv,

	"if": TOKIf, "else": TOKElse, "while": TOKWhile,
	"for": TOKFor, "match": TOKMatch, "defer": TOKDefer,

	"return": TOKReturn, "continue": TOKContinue, "break": TOKBreak,

	"const": TOKConst, "struct": TOKStruct,
	"enum": TOKEnum, "type": TOKType,
}

func StringIsKeyword(str string) bool {
	_, ok := keywords[str]
	return ok
}

func KeywordToKind(str string) (TokenKind, error) {
	kind, ok := keywords[str]
	if !ok {
		return 0, fmt.Errorf("%s is not a valid keyword", str)
	}
	return kind, nil
}
