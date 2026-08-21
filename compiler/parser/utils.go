package parser

import (
	"fmt"
	"naoslang/tokenizer/tokens"
)

func (p *parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.tokens[p.current].Kind == tokens.TOKEof
}

func (p *parser) isAtEndNToken(n int) bool {
	ptr := p.current + n
	return ptr >= len(p.tokens) || p.tokens[ptr].Kind == tokens.TOKEof
}

func (p *parser) readToken() tokens.Token {
	if p.isAtEnd() {
		return p.tokens[len(p.tokens)-1] // should be always eof
	}
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *parser) peekToken() tokens.Token {
	if p.isAtEnd() {
		return p.tokens[len(p.tokens)-1] // should be always eof
	}
	return p.tokens[p.current]
}

func (p *parser) peekNToken(n int) tokens.Token {
	if p.isAtEndNToken(n) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current+n]
}

func (p *parser) expect(kind tokens.TokenKind, errorMsg string) (tokens.Token, error) {
	t := p.peekToken()
	if t.Kind == kind {
		return p.readToken(), nil
	}
	return tokens.Token{}, fmt.Errorf("syntactic error at the position %v: %s (found insted %v)", t.Start, errorMsg, t.Kind)
}

func (p *parser) peekPrecedence() tokens.Prec {
	return p.peekToken().Kind.Precedence()
}

func (p *parser) cleanupTokens() error {
	ntokens := []tokens.Token{}
	for _, tok := range p.tokens {
		if tok.Kind.IsError() {
			return fmt.Errorf("syntax: %s", tok.Lexeme)
		}
		if !tok.Kind.IsUseless() {
			ntokens = append(ntokens, tok)
		}

		switch tok.Kind {
		case tokens.TOKIncrement:
			ntokens = append(ntokens, tokens.Token{
				Kind:   tokens.TOKPlusAssign,
				Lexeme: "",
				Start:  tok.Start,
				End:    tok.End,
			})
			ntokens = append(ntokens, tokens.Token{
				Kind:   tokens.TOKInt,
				Lexeme: "1",
				Start:  tok.Start,
				End:    tok.End,
			})
			continue

		case tokens.TOKDecrement:
			ntokens = append(ntokens, tokens.Token{
				Kind:   tokens.TOKMinusAssign,
				Lexeme: "",
				Start:  tok.Start,
				End:    tok.End,
			})
			ntokens = append(ntokens, tokens.Token{
				Kind:   tokens.TOKInt,
				Lexeme: "1",
				Start:  tok.Start,
				End:    tok.End,
			})
			continue
		}
	}

	p.tokens = ntokens
	return nil
}
