package tokens

import (
	"fmt"
	"strings"
)

func (t TokenKind) IsValid() bool           { return tok_start < t && t < tok_count }
func (t TokenKind) IsParentheses() bool     { return tok_paren_beg < t && t < tok_paren_end }
func (t TokenKind) IsSpecial() bool         { return tok_special_beg < t && t < tok_special_end }
func (t TokenKind) IsComment() bool         { return t == TOKMonoComment || t == TOKMultiComment }
func (t TokenKind) IsPunctuation() bool     { return tok_punct_beg < t && t < tok_punct_end }
func (t TokenKind) IsOperator() bool        { return tok_op_beg < t && t < tok_op_end }
func (t TokenKind) IsNormalOperator() bool  { return tok_nop_beg < t && t < tok_nop_beg }
func (t TokenKind) IsBooleanOperator() bool { return tok_bop_beg < t && t < tok_bop_end }
func (t TokenKind) IsBitwiseOperator() bool { return tok_bwop_beg < t && t < tok_bwop_beg }
func (t TokenKind) IsSpecialOperator() bool { return tok_sop_beg < t && t < tok_sop_beg }
func (t TokenKind) IsLiteral() bool         { return tok_lit_beg < t && t < tok_lit_end }
func (t TokenKind) IsKeyword() bool         { return tok_keyword_beg < t && t < tok_keyword_end }
func (t TokenKind) IsUseless() bool         { return t == TOKMonoComment || t == TOKMultiComment }
func (t TokenKind) IsError() bool           { return t == TOKInvalid || t == TOKNotImplemented }

var assign = map[TokenKind]struct{}{
	TOKAssign: {}, TOKPlusAssign: {}, TOKMinusAssign: {},
	TOKStarAssign: {}, TOKSlashAssign: {}, TOKPercentAssign: {},
	TOKBwAndAssign: {}, TOKBwOrAssign: {}, TOKBwXorAssign: {},
	TOKBwLShiftAssign: {}, TOKBwRShiftAssign: {},
}

func (t TokenKind) IsAssign() bool {
	_, ok := assign[t]
	return ok
}

// +------------+
// | Precedence |
// +------------+

type Prec int8

const (
	PrecedenceLowest     Prec = iota
	PrecedenceOr              // ||
	PrecedenceAnd             // &&
	PrecedenceBwOr            // |
	PrecedenceBwXor           // ^
	PrecedenceBwAnd           // &
	PrecedenceComparison      // == , != , < , > , <= , >=
	PrecedenceBwShift         // << , >>
	PrecedenceSum             // + , -
	PrecedenceProduct         // * , / , %
	PrecedenceUnary           // ! , - , ~
	PrecedenceCall            // function() or istance.method()
)

var precedence = map[TokenKind]Prec{
	TOKOr: PrecedenceOr, TOKAnd: PrecedenceAnd,
	TOKBwOr: PrecedenceBwOr, TOKBwXor: PrecedenceBwXor, TOKBwAnd: PrecedenceBwAnd,
	TOKEqual: PrecedenceComparison, TOKNotEqual: PrecedenceComparison, TOKLess: PrecedenceComparison,
	TOKLessEqual: PrecedenceComparison, TOKGreater: PrecedenceComparison, TOKGreaterEqual: PrecedenceComparison,
	TOKBwLShift: PrecedenceBwShift, TOKBwRShift: PrecedenceBwShift,
	TOKPlus: PrecedenceSum, TOKMinus: PrecedenceSum,
	TOKStar: PrecedenceProduct, TOKSlash: PrecedenceProduct, TOKPercent: PrecedenceProduct,
	TOKLParen: PrecedenceCall, TOKPeriod: PrecedenceCall,
}

func (t TokenKind) Precedence() Prec {
	prec, ok := precedence[t]
	if ok {
		return prec
	}
	return PrecedenceLowest
}
func (t TokenKind) IsUnary() bool { return t == TOKMinus || t == TOKNot || t == TOKBwNot }

// +------+
// | Misc |
// +------+

func PrintTokens(toks []Token, showPos bool) {
	headers := []string{"TOKEN", "LEXEME"}
	if showPos {
		headers = append([]string{"POSITION"}, headers...)
	}

	var rows [][]string
	for _, tok := range toks {
		lexeme := fmt.Sprintf("%v", tok.Lexeme)
		if lexeme == "" {
			lexeme = "-"
		}
		kindStr := fmt.Sprintf("%v", tok.Kind)

		var row []string
		if showPos {
			posStr := fmt.Sprintf("%v - %v", tok.Start, tok.End)
			row = []string{posStr, kindStr, lexeme}
		} else {
			row = []string{kindStr, lexeme}
		}
		rows = append(rows, row)
	}

	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, val := range row {
			if len(val) > colWidths[i] {
				colWidths[i] = len(val)
			}
		}
	}

	var fmtParts []string
	for _, w := range colWidths {
		fmtParts = append(fmtParts, fmt.Sprintf("%%-%ds", w))
	}
	rowFmt := "| " + strings.Join(fmtParts, " | ") + " |\n"

	totalWidth := len(headers)*3 + 1
	for _, w := range colWidths {
		totalWidth += w
	}
	divider := strings.Repeat("-", totalWidth)

	fmt.Println(divider)
	fmt.Printf(rowFmt, convertToInterfaceSlice(headers)...)
	fmt.Println(divider)
	for _, row := range rows {
		fmt.Printf(rowFmt, convertToInterfaceSlice(row)...)
	}
	fmt.Println(divider)
}

func convertToInterfaceSlice(strs []string) []interface{} {
	ifs := make([]interface{}, len(strs))
	for i, s := range strs {
		ifs[i] = s
	}
	return ifs
}
