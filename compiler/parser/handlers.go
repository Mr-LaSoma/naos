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
	_, err := p.expect(tokens.TOKPackage, "expected 'package' declaration at the start of the file")
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

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after type declaration")
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

			_, err := p.expect(tokens.TOKRBrackets, "expected ']' after '*' in multi-pointer declaration")
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

		_, err = p.expect(tokens.TOKRBraces, "expected '}' at the end of struct declaration")
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

		_, err = p.expect(tokens.TOKSemiColon, "expected ';' after field declaration")
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

// +-----------+
// | Variables |
// +-----------+

func (p *parser) handleVariableDecl(isPublic bool, isConst bool) (ast.ASTNode, error) {
	s := "variable"
	if isConst {
		s = "constant"
	}
	nameTok, err := p.expect(tokens.TOKId, fmt.Sprintf("expected %s name", s))
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKColon, fmt.Sprintf("expected ':' after %s name", s))
	if err != nil {
		return nil, err
	}

	var fieldType ast.ASTNode = nil
	if p.peekToken().Kind != tokens.TOKAssign {
		fieldType, err = p.handleTypeSignature()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(tokens.TOKAssign, fmt.Sprintf("expected '=' after type of %s", s))
	if err != nil {
		return nil, err
	}

	exprNode, err := p.handleExpressions(tokens.PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after declaration")
	if err != nil {
		return nil, err
	}

	return &ast.VariableDeclNode{
		IsPublic:   isPublic,
		IsConst:    isConst,
		Name:       nameTok.Lexeme,
		Type:       fieldType,
		Expression: exprNode,
	}, nil
}

// +-------------+
// | Expressions |
// +-------------+

func (p *parser) handleExpressions(precedence tokens.Prec) (ast.ASTNode, error) {
	leftNode, err := p.handleExpressionPrefix()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && precedence < p.peekPrecedence() && p.peekToken().Kind != tokens.TOKSemiColon {
		leftNode, err = p.handleExpressionInfix(leftNode)
		if err != nil {
			return nil, err
		}
	}

	return leftNode, nil
}

func (p *parser) handleExpressionPrefix() (ast.ASTNode, error) {
	tok := p.peekToken()

	switch tok.Kind {
	case tokens.TOKInt:
		p.readToken()
		return &ast.IntLiteralNode{Value: tok.Lexeme}, nil
	case tokens.TOKFloat:
		p.readToken()
		return &ast.FloatLiteralNode{Value: tok.Lexeme}, nil
	case tokens.TOKString:
		p.readToken()
		return &ast.StringLiteralNode{Value: tok.Lexeme}, nil
	case tokens.TOKChar:
		p.readToken()
		return &ast.CharLiteralNode{Value: tok.Lexeme}, nil
	case tokens.TOKId:
		p.readToken()
		return &ast.IdentifierRefNode{Name: tok.Lexeme}, nil
	case tokens.TOKAt:
		p.readToken()
		_, err := p.expect(tokens.TOKId, "expected action name")
		if err != nil {
			return nil, err
		}

		panic("not yet implemented '@' actions")

	case tokens.TOKLParen:
		p.readToken()
		expr, err := p.handleExpressions(tokens.PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKRParen, "expect ')' to close grouped expression")
		if err != nil {
			return nil, err
		}
		return expr, nil
	}

	if tok.Kind.IsUnary() {
		p.readToken()

		right, err := p.handleExpressions(tokens.PrecedenceUnary)
		if err != nil {
			return nil, err
		}

		return &ast.PrefixExpressionNode{Operator: tok.Kind.String(), Right: right}, nil
	}

	return nil, fmt.Errorf("expected start of expression, found token %v (value %s)", tok.Kind, tok.Lexeme)
}

func (p *parser) handleExpressionInfix(leftNode ast.ASTNode) (ast.ASTNode, error) {
	tok := p.peekToken()

	// function call
	if tok.Kind == tokens.TOKLParen {
		p.readToken()

		callNode := &ast.CallExpressionNode{
			Callee:    leftNode,
			Arguments: []ast.ASTNode{},
		}

		isFirstArgument := true
		for p.peekToken().Kind != tokens.TOKRParen && !p.isAtEnd() {
			if !isFirstArgument {
				_, err := p.expect(tokens.TOKComma, "expected ',' between arguments of function")
				if err != nil {
					return nil, err
				}
			}

			arg, err := p.handleExpressions(tokens.PrecedenceLowest)
			if err != nil {
				return nil, err
			}

			callNode.Arguments = append(callNode.Arguments, arg)
			isFirstArgument = false
		}

		_, err := p.expect(tokens.TOKRParen, "expected ')' after function arguments")
		if err != nil {
			return nil, err
		}

		return callNode, nil
	}

	currentPrec := p.peekPrecedence()
	p.readToken()

	rightNode, err := p.handleExpressions(currentPrec)
	if err != nil {
		return nil, err
	}

	return &ast.BinaryExpressionNode{
		Left:     leftNode,
		Operator: tok.Kind.String(),
		Right:    rightNode,
	}, nil
}

func (p *parser) handleAssignement(tok tokens.Token) (ast.ASTNode, error) {
	opTok := p.peekToken()
	if !opTok.Kind.IsAssign() {
		return nil, fmt.Errorf("expected assignement operator after identifier, found %s", opTok.Kind.String())
	}
	p.readToken()

	rightExpr, err := p.handleExpressions(tokens.PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after assignement statement")
	if err != nil {
		return nil, err
	}

	return &ast.AssignementNode{
		Left:     tok.Lexeme,
		Operator: opTok.Kind.String(),
		Right:    rightExpr,
	}, nil
}

// +-----------+
// | Functions |
// +-----------+

func (p *parser) handleFunctionDecl(isPublic bool) (ast.ASTNode, error) {
	nameTok, err := p.expect(tokens.TOKId, "expected function name")
	if err != nil {
		return nil, err
	}

	params, err := p.handleParameters()
	if err != nil {
		return nil, err
	}

	var returnType ast.ASTNode = nil
	if p.peekToken().Kind == tokens.TOKArrow {
		p.readToken()
		returnType, err = p.handleTypeSignature()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FuncDeclNode{
		IsPublic:   isPublic,
		Name:       nameTok.Lexeme,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

type tempParam struct {
	IsConst bool
	Name    string
}

func (p *parser) handleParameters() ([]ast.Parameter, error) {
	params := []ast.Parameter{}

	_, err := p.expect(tokens.TOKLParen, "expected '(' after function name")
	if err != nil {
		return nil, err
	}

	isFirstParam := true

	tempGroup := []tempParam{}

	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRParen {
		if !isFirstParam {
			_, err = p.expect(tokens.TOKComma, "expected ',' between function parameters")
			if err != nil {
				return nil, err
			}
		}
		isFirstParam = false

		isConst := false
		if p.peekToken().Kind == tokens.TOKConst {
			p.readToken()
			isConst = true
		}

		paramNameTok, err := p.expect(tokens.TOKId, "expected parameter name")
		if err != nil {
			return nil, err
		}

		tempGroup = append(tempGroup, tempParam{
			IsConst: isConst,
			Name:    paramNameTok.Lexeme,
		})

		if p.peekToken().Kind == tokens.TOKColon {
			p.readToken()

			sharedType, err := p.handleTypeSignature()
			if err != nil {
				return nil, err
			}

			for _, tp := range tempGroup {
				params = append(params, ast.Parameter{
					IsConst: tp.IsConst,
					Name:    tp.Name,
					Type:    sharedType,
				})
			}

			tempGroup = []tempParam{}
		}
	}

	if len(tempGroup) > 0 {
		return nil, fmt.Errorf("expected ':' at the end of the parameters")
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after parameters")
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (p *parser) handleBlock() ([]ast.ASTNode, error) {
	body := []ast.ASTNode{}

	_, err := p.expect(tokens.TOKLBraces, "expected '{' after function name")
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRBraces {
		var statem ast.ASTNode
		tok := p.peekToken()

		switch tok.Kind {
		case tokens.TOKLet, tokens.TOKConst:
			p.readToken()
			isConst := tok.Kind == tokens.TOKConst
			statem, err = p.handleVariableDecl(false, isConst)

		case tokens.TOKId:
			next := p.peekNToken(1)

			if next.Kind.IsAssign() {
				p.readToken()
				statem, err = p.handleAssignement(tok)
			} else {
				statem, err = p.handleExpressions(tokens.PrecedenceLowest)
				if err != nil {
					_, err = p.expect(tokens.TOKSemiColon, "expected ';' after expression")
				}
			}

		case tokens.TOKAt:
			p.readToken()
			_, err = p.expect(tokens.TOKId, "expected action name")
			if err != nil {
				return nil, err
			}

			panic("not yet implemented '@' actions")

		default:
			statem, err = p.handleExpressions(tokens.PrecedenceLowest)
			if err != nil {
				_, err = p.expect(tokens.TOKSemiColon, "expected ';' after expression")
			}
		}

		if err != nil {
			return nil, err
		}
		body = append(body, statem)
	}

	_, err = p.expect(tokens.TOKRBraces, "expected '}' after function body")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// +---------------+
// | TypeExtension |
// +---------------+

func (p *parser) HandleTypeExtension(targetType string) (ast.ASTNode, error) {
	_, err := p.expect(tokens.TOKLBraces, "expected '{' after extension operator")
	if err != nil {
		return nil, err
	}

	extNode := &ast.TypeExtensionNode{
		TargetType: targetType,
		Methods:    []ast.ASTNode{},
	}

	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRBraces {
		var methodNode ast.ASTNode
		tok := p.peekToken()

		switch tok.Kind {
		case tokens.TOKPub, tokens.TOKPriv:
			p.readToken()
			isPublic := tok.Kind == tokens.TOKPub

			methodNode, err = p.handleFunctionDecl(isPublic)

		case tokens.TOKId:
			methodNode, err = p.handleFunctionDecl(false)

		case tokens.TOKAt:
			panic("not yet implemented the '@' symbol in the type extension")

		default:
			return nil, fmt.Errorf("unexpected token %v (value %v) in type %s extension block", tok.Kind, tok.Lexeme, targetType)
		}

		if err != nil {
			return nil, err
		}
		extNode.Methods = append(extNode.Methods, methodNode)
	}

	_, err = p.expect(tokens.TOKRBraces, "expected '}' after type extension block")
	if err != nil {
		return nil, err
	}

	return extNode, nil
}
