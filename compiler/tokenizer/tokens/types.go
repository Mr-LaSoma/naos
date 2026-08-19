package tokens

import "strconv"

type TokenKind int8

const (
	tok_start TokenKind = iota

	// +-----------------------+
	// | start of the specials |
	// +---------------- ------+
	tok_special_beg
	TOKNotImplemented // token the lexer didn'y yet implement
	TOKInvalid        // invalid syntaxt found
	TOKEof            // EOF
	TOKMonoComment    // //
	TOKMultiComment   // /*

	// end of the parentheses
	tok_special_end

	// +--------------------------+
	// | start of the parentheses |
	// +--------------------------+
	tok_paren_beg
	TOKLParen    // (
	TOKRParen    // )
	TOKLBrackets // [
	TOKRBrackets // ]
	TOKLBraces   // {
	TOKRBraces   // }

	// end of the parentheses
	tok_paren_end

	tok_count
)

var tokens = [...]string{
	TOKNotImplemented: "NOT IMPLEMENTED",
	TOKInvalid:        "invalid",
	TOKEof:            "EOF",
	TOKMonoComment:    "//",
	TOKMultiComment:   "/*",

	TOKLParen:    "(",
	TOKRParen:    ")",
	TOKLBrackets: "[",
	TOKRBrackets: "]",
	TOKLBraces:   "{",
	TOKRBraces:   "}",
}

func (t TokenKind) String() string {
	s := ""
	if tok_start <= t && t <= tok_count {
		s = tokens[t]
	}
	if s == "" {
		return "unknown(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}
