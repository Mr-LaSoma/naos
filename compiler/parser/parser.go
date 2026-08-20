package parser

import (
	"naoslang/parser/ast"
	"naoslang/tokenizer/tokens"
)

type parser struct {
	tokens  []tokens.Token
	current int
}

func New(tokens []tokens.Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() (*ast.ASTFile, error) {
	err := p.dropUselessAndErrors()
	if err != nil {
		return nil, err
	}

	fileAst := &ast.ASTFile{
		Imports: []ast.ASTNode{},
		Decls:   []ast.ASTNode{},
	}
	err = p.handlePackageDecl(fileAst)
	if err != nil {
		return nil, err
	}

	err = p.handleImports(fileAst)
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() {
		tok := p.peekToken()
		switch tok.Kind {
		case tokens.TOKType:
			p.readToken()
			typeNode, err := p.handleTypeDecl(false)
			if err != nil {
				return nil, err
			}
			fileAst.Decls = append(fileAst.Decls, typeNode)
			continue

		case tokens.TOKPub, tokens.TOKPriv:
			p.readToken()
			isPublic := tok.Kind == tokens.TOKPub
			ntok := p.peekToken()

			switch ntok.Kind {
			case tokens.TOKType:
				p.readToken()
				typeNode, err := p.handleTypeDecl(isPublic)
				if err != nil {
					return nil, err
				}
				fileAst.Decls = append(fileAst.Decls, typeNode)
				continue
			}
		}

		p.readToken()
	}

	return fileAst, nil
}
