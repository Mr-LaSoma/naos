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

var operators = map[string]TokenKind{
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

const operatorsStart = "+-*/%=!&|<>^~#@:"

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

	"let": TOKLet, "const": TOKConst, "struct": TOKStruct,
	"enum": TOKEnum, "type": TOKType, "interface": TOKInterface,
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

// +---------+
// | Numbers |
// +---------+

type numericBase int8

const (
	NumericBaseBinary numericBase = iota
	NumericBaseOctal
	NumericBaseDecimal
	NumericBaseHexadecimal
)

func RuneToNumericBase(ch rune) (numericBase, error) {
	switch ch {
	case 'b', 'B':
		return NumericBaseBinary, nil
	case 'o', 'O':
		return NumericBaseOctal, nil
	case 'x', 'X':
		return NumericBaseHexadecimal, nil
	}
	return 0, fmt.Errorf("%c is not a valid numeric base", ch)
}

func IsBinaryDigit(ch rune) bool  { return ch == '0' || ch == '1' }
func IsOctalDigit(ch rune) bool   { return '0' <= ch && ch <= '7' }
func IsDecimalDigit(ch rune) bool { return '0' <= ch && ch <= '9' }
func IsHexDigit(ch rune) bool {
	return IsDecimalDigit(ch) || ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func IsValidBaseDigit(ch rune, base numericBase) bool {
	switch base {
	case NumericBaseBinary:
		return IsBinaryDigit(ch)
	case NumericBaseOctal:
		return IsOctalDigit(ch)
	case NumericBaseDecimal:
		return IsDecimalDigit(ch)
	case NumericBaseHexadecimal:
		return IsHexDigit(ch)
	}
	return false
}

func (b numericBase) String() string {
	switch b {
	case NumericBaseBinary:
		return "binary"
	case NumericBaseOctal:
		return "octal"
	case NumericBaseDecimal:
		return "decimal"
	case NumericBaseHexadecimal:
		return "hexadecimal"
	}
	return "UNKNOWN"
}

// +---------+
// | Strings |
// +---------+

const baseEscapeSeq = `abfnrtv\'"`

func IsBaseEscapeSeq(ch rune) bool {
	for _, opBase := range baseEscapeSeq {
		if ch == opBase {
			return true
		}
	}
	return false
}

type escapeKind int8

const (
	EscapeKindOctal escapeKind = iota
	EscapeKindHexadecimal
	EscapeKindUnicode16
	EscapeKindUnicode32
)

func (e escapeKind) NumbersNeeded() int {
	switch e {
	case EscapeKindOctal:
		return 3
	case EscapeKindHexadecimal:
		return 2
	case EscapeKindUnicode16:
		return 4
	case EscapeKindUnicode32:
		return 8
	}
	return 0
}

func (e escapeKind) ToNumericBase() numericBase {
	switch e {
	case EscapeKindOctal:
		return NumericBaseOctal
	case EscapeKindHexadecimal, EscapeKindUnicode16, EscapeKindUnicode32:
		return NumericBaseHexadecimal
	}
	return 0
}

func (e escapeKind) String() string {
	switch e {
	case EscapeKindOctal:
		return "octal"
	case EscapeKindHexadecimal:
		return "hexadecimal"
	case EscapeKindUnicode16:
		return "unicode16"
	case EscapeKindUnicode32:
		return "unicode32"
	}
	return "UNKNOWN"
}

func RuneToEscapeKind(ch rune) (escapeKind, error) {
	switch ch {
	case 'x':
		return EscapeKindHexadecimal, nil
	case 'u':
		return EscapeKindUnicode16, nil
	case 'U':
		return EscapeKindUnicode32, nil
	}
	return 0, fmt.Errorf("%c is not a valid escape sequence base", ch)
}

// +---------+
// | Actions |
// +---------+

func StringIsActionWParen(str string) bool {
	switch str {
	case "null", "unreachable":
		return false
	}
	return true
}
