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
	err := p.cleanupTokens() // this one is preatty sus
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

		case tokens.TOKAt:
			p.readToken()
			identTok, errt := p.expect(tokens.TOKId, "expected compiler action name")
			if errt != nil {
				return nil, errt
			}

			switch identTok.Lexeme {
			case "inline", "noinline", "extern":
				attrs, errt := p.handleFunctionAttributes()
				if errt != nil {
					return nil, errt
				}

				ntok := p.peekToken()
				isPublic := false
				if ntok.Kind == tokens.TOKPub || ntok.Kind == tokens.TOKPriv {
					p.readToken()
					isPublic = ntok.Kind == tokens.TOKPub
				}

				astNode, errt := p.handleFunctionDecl(isPublic)
				if errt != nil {
					return nil, errt
				}

				methodNode, ok := astNode.(*ast.FuncDeclNode)
				if !ok {
					panic("handle function declaration doesn't returns a FuncDeclNode")
				}
				methodNode.Attributes = attrs
				fileAst.Decls = append(fileAst.Decls, methodNode)

			default:
				panic("not yet implemented the '@' action in global decl")
			}
		}

		p.readToken()
	}

	return fileAst, nil
}
