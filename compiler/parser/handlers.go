package parser

import (
	"fmt"
	"naoslang/parser/ast"
	"naoslang/tokenizer/tokens"
)

// +----------+
// | Packages |
// +----------+

func (p *parser) handlePackageDecl(fileAst *ast.ASTFile) error {
	_, err := p.expect(tokens.TOKPackage, "expected 'package' decleration at the start of the file")
	if err != nil {
		return err
	}
	nameTok, err := p.expect(tokens.TOKId, "expected package name after 'package'")
	if err != nil {
		return err
	}
	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after package name")
	if err != nil {
		return err
	}

	fileAst.PackageName = nameTok.Lexeme
	return nil
}

// +---------+
// | Imports |
// +---------+

func (p *parser) handleImports(fileAst *ast.ASTFile) error {
	for !p.isAtEnd() {
		tok := p.peekToken()
		if tok.Kind == tokens.TOKUsing {
			p.readToken()
			importNode, err := p.handleGlobalImport()
			if err != nil {
				return err
			}
			fileAst.Imports = append(fileAst.Imports, importNode)
		} else if p.isAliasImport() {
			// can skip id, ::, @ and import
			aliasTok := p.readToken()
			p.readToken()
			p.readToken()
			p.readToken()

			importNode, err := p.handleAliasImport(aliasTok.Lexeme)
			if err != nil {
				return err
			}
			fileAst.Imports = append(fileAst.Imports, importNode)
		} else {
			break
		}
	}
	return nil
}

func (p *parser) handleGlobalImport() (ast.ASTNode, error) {
	_, err := p.expect(tokens.TOKAt, "expected '@import' action after 'using'")
	if err != nil {
		return nil, err
	}

	actionTok, err := p.expect(tokens.TOKId, "expected 'import' identifier after '@'")
	if err != nil || actionTok.Lexeme != "import" {
		return nil, err
	}

	_, err = p.expect(tokens.TOKLParen, "expected '(' after '@import'")
	if err != nil {
		return nil, err
	}

	pathTok, err := p.expect(tokens.TOKString, "expected string literal inside '@import(...)'")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after import path")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after import statement")
	if err != nil {
		return nil, err
	}

	cleanPath := pathTok.Lexeme[1 : len(pathTok.Lexeme)-1]
	return &ast.ImportGlobalNode{Path: cleanPath}, nil
}

func (p *parser) handleAliasImport(alias string) (ast.ASTNode, error) {
	_, err := p.expect(tokens.TOKLParen, "expected '(' after '@import'")
	if err != nil {
		return nil, err
	}

	pathTok, err := p.expect(tokens.TOKString, "expected string literal inside '@import(...)'")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after import path")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after import statement")
	if err != nil {
		return nil, err
	}

	cleanPath := pathTok.Lexeme[1 : len(pathTok.Lexeme)-1]
	return &ast.ImportAliasNode{Alias: alias, Path: cleanPath}, nil
}

func (p *parser) isAliasImport() bool {
	if p.isAtEnd() {
		return false
	}

	tok := p.peekToken()
	importTok := p.peekNToken(3) // id :: @ id <- the 4° one
	return tok.Kind == tokens.TOKId && p.peekNToken(1).Kind == tokens.TOKDoubleColon &&
		p.peekNToken(2).Kind == tokens.TOKAt && (importTok.Kind == tokens.TOKId && importTok.Lexeme == "import")
}

// +-------+
// | Types |
// +-------+

func (p *parser) handleTypeDecl(isPublic bool) (ast.ASTNode, error) {
	nameTok, err := p.expect(tokens.TOKId, "expected type name after 'type'")
	if err != nil {
		return nil, err
	}

	declTok := p.peekToken()
	var isNewType bool

	if declTok.Kind == tokens.TOKAssign || declTok.Kind == tokens.TOKDoubleColon {
		p.readToken()
		isNewType = declTok.Kind == tokens.TOKDoubleColon
	} else {
		return nil, fmt.Errorf("expected '=' or ':' after type name '%s', found '%v'", nameTok.Lexeme, declTok.Kind)
	}

	typeSignature, err := p.handleTypeSignature()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after type decleration")
	if err != nil {
		return nil, err
	}

	if isNewType {
		return &ast.NewTypeNode{
			IsPublic: isPublic,
			Name:     nameTok.Lexeme,
			BaseType: typeSignature,
		}, nil
	}
	return &ast.TypeAliasNode{
		IsPublic:   isPublic,
		Name:       nameTok.Lexeme,
		SourceType: typeSignature,
	}, nil
}

func (p *parser) handleTypeSignature() (ast.ASTNode, error) {
	tok := p.peekToken()
	switch tok.Kind {
	case tokens.TOKLBrackets:
		p.readToken()
		ntok := p.peekToken()

		if ntok.Kind == tokens.TOKStar {
			p.readToken()

			_, err := p.expect(tokens.TOKRBrackets, "expected ']' after '*' in multi-pointer decleration")
			if err != nil {
				return nil, err
			}
			elemType, err := p.handleTypeSignature()
			if err != nil {
				return nil, err
			}

			return &ast.MultiPointerTypeNode{ElementType: elemType}, nil
		}

		return nil, fmt.Errorf("not yet implemented the []T and [N]T types")

	case tokens.TOKStar:
		p.readToken()
		elemType, err := p.handleTypeSignature()
		if err != nil {
			return nil, err
		}

		return &ast.PointerTypeNode{ElementType: elemType}, nil

	case tokens.TOKStruct:
		p.readToken()
		_, err := p.expect(tokens.TOKLBraces, "expected '{' after 'struct'")
		if err != nil {
			return nil, err
		}

		structNode, err := p.handleStructFields()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKRBraces, "expected '}' at the end of struct decleration")
		if err != nil {
			return nil, err
		}
		return structNode, nil

	case tokens.TOKId:
		p.readToken()
		return &ast.TypeReferanceNode{Name: tok.Lexeme}, nil
	}

	return nil, fmt.Errorf("expected valid type signature, found %v", tok.Kind)
}

func (p *parser) handleStructFields() (ast.ASTNode, error) {
	structNode := &ast.StructNode{Fields: []ast.StructField{}}
	tok := p.peekToken()

	for !p.isAtEnd() && tok.Kind != tokens.TOKRBraces {
		ntok := p.peekToken()
		isPublic := false

		if ntok.Kind == tokens.TOKPub || ntok.Kind == tokens.TOKPriv {
			p.readToken()
			isPublic = ntok.Kind == tokens.TOKPub
		}

		fieldNameTok, err := p.expect(tokens.TOKId, "expected field name inside struct")
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKColon, "expected ':' after field name")
		if err != nil {
			return nil, err
		}

		fieldType, err := p.handleTypeSignature()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKSemiColon, "expected ';' after field decleration")
		if err != nil {
			return nil, err
		}

		structNode.Fields = append(structNode.Fields, ast.StructField{
			IsPublic: isPublic,
			Name:     fieldNameTok.Lexeme,
			Type:     fieldType,
		})
		tok = p.peekToken()
	}
	return structNode, nil
}
