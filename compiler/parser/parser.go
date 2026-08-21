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
	err := p.cleanupTokens()
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
		case tokens.TOKId:
			ntok := p.peekNToken(1)
			if ntok.Kind == tokens.TOKLParen {
				declNode, err := p.handleFunctionDecl(false)
				if err != nil {
					return nil, err
				}
				fileAst.Decls = append(fileAst.Decls, declNode)
			}

			if ntok.Kind == tokens.TOKHashtag {
				p.readToken()
				p.readToken()

				extNode, err := p.HandleTypeExtension(tok.Lexeme)
				if err != nil {
					return nil, err
				}
				fileAst.Decls = append(fileAst.Decls, extNode)
			}

			continue

		case tokens.TOKLet, tokens.TOKConst:
			p.readToken()
			isConst := tok.Kind == tokens.TOKConst
			declNode, err := p.handleVariableDecl(false, isConst)
			if err != nil {
				return nil, err
			}
			fileAst.Decls = append(fileAst.Decls, declNode)
			continue

		case tokens.TOKPub, tokens.TOKPriv:
			p.readToken()
			isPublic := tok.Kind == tokens.TOKPub
			ntok := p.peekToken()

			switch ntok.Kind {
			case tokens.TOKId:
				if p.peekNToken(1).Kind == tokens.TOKLParen {
					declNode, err := p.handleFunctionDecl(isPublic)
					if err != nil {
						return nil, err
					}
					fileAst.Decls = append(fileAst.Decls, declNode)
				}
				continue

			case tokens.TOKLet, tokens.TOKConst:
				p.readToken()
				isConst := tok.Kind == tokens.TOKConst
				declNode, err := p.handleVariableDecl(isPublic, isConst)
				if err != nil {
					return nil, err
				}
				fileAst.Decls = append(fileAst.Decls, declNode)
				continue

			case tokens.TOKType:
				p.readToken()
				typeNode, err := p.handleTypeDecl(isPublic)
				if err != nil {
					return nil, err
				}
				fileAst.Decls = append(fileAst.Decls, typeNode)
				continue
			}
		case tokens.TOKType:
			p.readToken()
			typeNode, err := p.handleTypeDecl(false)
			if err != nil {
				return nil, err
			}
			fileAst.Decls = append(fileAst.Decls, typeNode)
			continue

		}

		p.readToken()
	}

	return fileAst, nil
}
