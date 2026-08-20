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

	// +---------------------------+
	// | start of the punctuations |
	// +---------------------------+
	tok_punct_beg
	TOKSemiColon // ;
	TOKComma     // ,
	TOKPeriod    // .

	// end of the puncts
	tok_punct_end

	// +------------------------+
	// | start of the operators |
	// +------------------------+
	tok_op_beg
	tok_nop_beg      // start of the normal ones
	TOKPlus          // +
	TOKMinus         // -
	TOKStar          // *
	TOKSlash         // /
	TOKPercent       // %
	TOKAssign        // =
	TOKPlusAssign    // +=
	TOKMinusAssign   // -=
	TOKStarAssign    // *=
	TOKSlashAssign   // /=
	TOKPercentAssign // %=
	TOKIncrement     // ++
	TOKDecrement     // --

	tok_nop_end // end of the normal ones

	tok_bop_beg     // start of the boolean ones
	TOKNot          // !
	TOKAnd          // &&
	TOKOr           // ||
	TOKLess         // <
	TOKGreater      // >
	TOKEqual        // ==
	TOKNotEqual     // !=
	TOKLessEqual    // <=
	TOKGreaterEqual // >=

	tok_bop_end // end of the boolean ones

	tok_bwop_beg      // start of the bitwise ones
	TOKBwAnd          // &
	TOKBwOr           // |
	TOKBwXor          // ^
	TOKBwNot          // ~
	TOKBwLShift       // <<
	TOKBwRShift       // >>
	TOKBwAndAssign    // &=
	TOKBwOrAssign     // |=
	TOKBwXorAssign    // ^=
	TOKBwLShiftAssign // <<=
	TOKBwRShiftAssign // >>=

	tok_bwop_end // end of the bitwise ones

	tok_sop_beg    // start of the special ones
	TOKHashtag     // #
	TOKAt          // @
	TOKArrow       // ->
	TOKBigArrow    // =>
	TOKColon       // :
	TOKDoubleColon // ::

	tok_sop_end // end of the special ones

	// end of the operators
	tok_op_end

	// +-------------------+
	// | start of literals |
	// +-------------------+
	tok_lit_beg
	TOKId

	// end of literals
	tok_lit_end

	// +-------------------+
	// | start of keywords |
	// +-------------------+
	tok_keyword_beg
	TOKPackage // package
	TOKUsing   // using
	TOKPub     // pub
	TOKPriv    // priv

	TOKIf    // if
	TOKElse  // else
	TOKWhile // while
	TOKFor   // for
	TOKMatch // match
	TOKDefer // defer

	TOKReturn   // return
	TOKContinue // continue
	TOKBreak    // break

	TOKConst  // const
	TOKStruct // struct
	TOKEnum   // enum
	TOKType   // type

	// end of keywords
	tok_keyword_end

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

	TOKColon:     ":",
	TOKSemiColon: ";",
	TOKComma:     ",",
	TOKPeriod:    ".",

	TOKPlus:          "+",
	TOKMinus:         "-",
	TOKStar:          "*",
	TOKSlash:         "/",
	TOKPercent:       "%",
	TOKAssign:        "=",
	TOKPlusAssign:    "+=",
	TOKMinusAssign:   "-=",
	TOKStarAssign:    "*=",
	TOKSlashAssign:   "/=",
	TOKPercentAssign: "%=",
	TOKIncrement:     "++",
	TOKDecrement:     "--",

	TOKNot:          "!",
	TOKAnd:          "&&",
	TOKOr:           "||",
	TOKLess:         "<",
	TOKGreater:      ">",
	TOKEqual:        "==",
	TOKNotEqual:     "!=",
	TOKLessEqual:    "<=",
	TOKGreaterEqual: ">=",

	TOKBwAnd:          "&",
	TOKBwOr:           "|",
	TOKBwXor:          "^",
	TOKBwNot:          "~",
	TOKBwLShift:       "<<",
	TOKBwRShift:       ">>",
	TOKBwAndAssign:    "&=",
	TOKBwOrAssign:     "|=",
	TOKBwXorAssign:    "^=",
	TOKBwLShiftAssign: "<<=",
	TOKBwRShiftAssign: ">>=",

	TOKHashtag:     "#",
	TOKAt:          "@",
	TOKArrow:       "->",
	TOKBigArrow:    "=>",
	TOKDoubleColon: "::",

	TOKId: "identifier",

	TOKPackage: "package",
	TOKUsing:   "using",
	TOKPub:     "pub",
	TOKPriv:    "priv",

	TOKIf:    "if",
	TOKElse:  "else",
	TOKWhile: "while",
	TOKFor:   "for",
	TOKMatch: "match",
	TOKDefer: "defer",

	TOKReturn:   "return",
	TOKContinue: "continue",
	TOKBreak:    "break",

	TOKConst:  "const",
	TOKStruct: "struct",
	TOKEnum:   "enum",
	TOKType:   "type",
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
