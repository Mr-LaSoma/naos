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

		} else if tok.Kind == tokens.TOKId && p.peekNToken(1).Kind == tokens.TOKDoubleColon {
			// can skip id, ::, @ and import
			aliasTok := p.readToken()

			_, err := p.expect(tokens.TOKDoubleColon, "expected '::' after import alias")
			if err != nil {
				return err
			}
			_, err = p.expect(tokens.TOKAt, "expected '@import' after ::")
			if err != nil {
				return err
			}
			importTok, err := p.expect(tokens.TOKId, "expected 'import' after '@'")
			if err != nil || importTok.Lexeme != "import" {
				return fmt.Errorf("expected 'import' after '@', found %s", importTok.Lexeme)
			}

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

	case tokens.TOKInterface:
		p.readToken()
		_, err := p.expect(tokens.TOKLBraces, "expected '{' after 'interface'")
		if err != nil {
			return nil, err
		}

		interfaceNode, err := p.handleInterfaceMethods()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKRBraces, "expected '}' at the end of interface declaration")
		if err != nil {
			return nil, err
		}
		return interfaceNode, nil

	case tokens.TOKEnum:
		p.readToken()
		_, err := p.expect(tokens.TOKLBraces, "expected '{' after 'enum'")
		if err != nil {
			return nil, err
		}

		enumNode, err := p.handleEnumVariants()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKRBraces, "expected '}' at the end of enum declaration")
		if err != nil {
			return nil, err
		}
		return enumNode, nil

	case tokens.TOKDollar:
		p.readToken()

		nameTok, err := p.expect(tokens.TOKId, "expected generic type name after '$'")
		if err != nil {
			return nil, err
		}

		var constraintNode ast.ASTNode = nil
		ntok := p.peekToken()
		if ntok.Kind.IsStartTypeSign() {
			constraintNode, err = p.handleTypeSignature()
			if err != nil {
				return nil, err
			}
		}

		return &ast.GenericTypeNode{
			Name:       nameTok.Lexeme,
			Constraint: constraintNode,
		}, nil

	case tokens.TOKId:
		p.readToken()
		var currentType ast.ASTNode = &ast.TypeReferanceNode{Name: tok.Lexeme}

		if p.peekToken().Kind == tokens.TOKLParen {
			p.readToken()

			instNode := &ast.GenericInstantiationNode{
				Left:         currentType,
				TypeArgument: []ast.ASTNode{},
			}

			isFirstArg := true
			for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRParen {
				if !isFirstArg {
					_, err := p.expect(tokens.TOKComma, "expected ',' between generic types arguments")
					if err != nil {
						return nil, err
					}
				}

				argType, err := p.handleTypeSignature()
				if err != nil {
					return nil, err
				}

				instNode.TypeArgument = append(instNode.TypeArgument, argType)
				isFirstArg = false
			}

			_, err := p.expect(tokens.TOKRParen, "expected ')' after generic arguments")
			if err != nil {
				return nil, err
			}

			currentType = instNode
		}

		return currentType, nil
	}

	return nil, fmt.Errorf("expected valid type signature, found %v", tok)
}

func (p *parser) handleEnumVariants() (ast.ASTNode, error) {
	enumNode := &ast.EnumLiteralNode{Variants: []ast.EnumVariant{}}
	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRBraces {
		variantNameTok, err := p.expect(tokens.TOKId, "expected variant name")
		if err != nil {
			return nil, err
		}

		variant := ast.EnumVariant{
			Name:   variantNameTok.Lexeme,
			Fields: []ast.StructField{},
		}

		if p.peekToken().Kind == tokens.TOKLParen {
			p.readToken()
			isFirstField := true

			for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRParen {
				if !isFirstField {
					_, err = p.expect(tokens.TOKComma, "expected ',' between variant fields")
					if err != nil {
						return nil, err
					}
				}

				var fieldName string
				var fieldType ast.ASTNode

				if p.peekToken().Kind == tokens.TOKId && p.peekNToken(1).Kind == tokens.TOKColon {
					fieldName = p.readToken().Lexeme
					p.readToken()

					fieldType, err = p.handleTypeSignature()
					if err != nil {
						return nil, err
					}
				} else {
					fieldType, err = p.handleTypeSignature()
					if err != nil {
						return nil, err
					}

					fieldName = fmt.Sprintf("item%d", len(variant.Fields))
				}

				isFirstField = false
				variant.Fields = append(variant.Fields, ast.StructField{
					IsPublic: true,
					Name:     fieldName,
					Type:     fieldType,
				})
			}

			_, err = p.expect(tokens.TOKRParen, "expected ')' after variant fields")
			if err != nil {
				return nil, err
			}

			_, err = p.expect(tokens.TOKComma, "expected ',' after variant declaration")
			if err != nil {
				return nil, err
			}
		}

		enumNode.Variants = append(enumNode.Variants, variant)
	}

	return enumNode, nil
}

func (p *parser) handleInterfaceMethods() (ast.ASTNode, error) {
	interfaceNode := &ast.InterfaceLiteralNode{Methods: []*ast.FuncSignatureNode{}}
	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRBraces {
		isMethodPublic := false
		if p.peekToken().Kind == tokens.TOKPub || p.peekToken().Kind == tokens.TOKPriv {
			isMethodPublic = p.readToken().Kind == tokens.TOKPub
		}

		methodName, err := p.expect(tokens.TOKId, "expected method name inside interface")
		if err != nil {
			return nil, err
		}

		params, err := p.handleParameters()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKArrow, "expected '->' after method parameters")
		if err != nil {
			return nil, err
		}

		returnType, err := p.handleTypeSignature()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokens.TOKComma, "expected ',' after method")
		if err != nil {
			return nil, err
		}

		interfaceNode.Methods = append(interfaceNode.Methods, &ast.FuncSignatureNode{
			IsPublic:   isMethodPublic,
			Name:       methodName.Lexeme,
			Parameter:  params,
			ReturnType: returnType,
		})
	}
	return interfaceNode, nil
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

		_, err = p.expect(tokens.TOKComma, "expected ',' after field declaration")
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
	case tokens.TOKIf:
		return p.handleIfStatement()
	case tokens.TOKWhile:
		return p.handleWhileStatement()
	case tokens.TOKFor:
		return p.handleForStatement()

	case tokens.TOKAt:
		p.readToken()
		actionTok, err := p.expect(tokens.TOKId, "expected action name")
		if err != nil {
			return nil, err
		}

		args := []ast.ASTNode{}

		if tokens.StringIsActionWParen(actionTok.Lexeme) {
			_, err := p.expect(tokens.TOKLParen, "expected '(' after action name")
			if err != nil {
				return nil, err
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

				args = append(args, arg)
				isFirstArgument = false
			}

			_, err = p.expect(tokens.TOKRParen, "expected ')' after function arguments")
			if err != nil {
				return nil, err
			}
		}

		return &ast.CompilerActionNode{
			Name:      actionTok.Lexeme,
			Left:      nil,
			Arguments: args,
		}, nil

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

	block, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (p *parser) handleExpressionInfix(leftNode ast.ASTNode) (ast.ASTNode, error) {
	tok := p.peekToken()

	switch tok.Kind {
	case tokens.TOKPeriod:
		p.readToken()

		memberTok, err := p.expect(tokens.TOKId, "expected member name after '.'")
		if err != nil {
			return nil, err
		}
		return &ast.MemberAccessNode{
			Left:   leftNode,
			Member: memberTok.Lexeme,
		}, nil

	case tokens.TOKLParen:
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

	case tokens.TOKAt:
		p.readToken()
		actionTok, err := p.expect(tokens.TOKId, "expected action name")
		if err != nil {
			return nil, err
		}

		args := []ast.ASTNode{}

		if tokens.StringIsActionWParen(actionTok.Lexeme) {
			_, err := p.expect(tokens.TOKLParen, "expected '(' after action name")
			if err != nil {
				return nil, err
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

				args = append(args, arg)
				isFirstArgument = false
			}

			_, err = p.expect(tokens.TOKRParen, "expected ')' after function arguments")
			if err != nil {
				return nil, err
			}
		}

		return &ast.CompilerActionNode{
			Name:      actionTok.Lexeme,
			Left:      nil,
			Arguments: args,
		}, nil
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

func (p *parser) handleFunctionAttributes() ([]ast.FuncAttribute, error) {
	attrs := []ast.FuncAttribute{}

	for p.peekToken().Kind == tokens.TOKAt {
		p.readToken()

		nameTok, err := p.expect(tokens.TOKId, "expected attribute name after '@'")
		if err != nil {
			return nil, err
		}

		attr := ast.FuncAttribute{
			Name: nameTok.Lexeme,
			Args: map[string]string{},
		}

		if tokens.StringIsAttrWParen(nameTok.Lexeme) {
			_, err := p.expect(tokens.TOKLParen, "expected '(' after attribute name")
			if err != nil {
				return nil, err
			}

			isFirstArg := true
			for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRParen {
				if !isFirstArg {
					_, err = p.expect(tokens.TOKComma, "expected ',' between attribute arguments")
					if err != nil {
						return nil, err
					}
				}

				keyTok, err := p.expect(tokens.TOKId, "expected argument key")
				if err != nil {
					return nil, err
				}

				if _, ok := attr.Args[keyTok.Lexeme]; ok {
					return nil, fmt.Errorf("expected unique argument key, duplicate found %v", keyTok.Lexeme)
				}

				_, err = p.expect(tokens.TOKAssign, "expected '=' after argument key")
				if err != nil {
					return nil, err
				}

				valTok, err := p.expect(tokens.TOKString, fmt.Sprintf("expected string value for key %v", keyTok.Lexeme))
				if err != nil {
					return nil, err
				}

				cleanVal := valTok.Lexeme[1 : len(valTok.Lexeme)-1]
				attr.Args[keyTok.Lexeme] = cleanVal

				isFirstArg = false
			}

			_, err = p.expect(tokens.TOKRParen, "expected ')' after attribute")
			if err != nil {
				return nil, err
			}
		}

		attrs = append(attrs, attr)
	}

	return attrs, nil
}

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

	fmt.Printf("TOKNEW: %v | %v\n", p.peekToken(), body)

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

			if p.peekToken().Kind.IsStartTypeSign() {
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
			} else {
				return nil, fmt.Errorf("expected valid type signature found %v", p.peekToken())
			}
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

func (p *parser) handleBlock() (ast.ASTNode, error) {
	_, err := p.expect(tokens.TOKLBraces, "expected '{' in block declaration")
	if err != nil {
		return nil, err
	}

	blockNode := &ast.BlockExpressionNode{
		Statements: []ast.ASTNode{},
		Value:      nil,
	}

	for !p.isAtEnd() && p.peekToken().Kind != tokens.TOKRBraces {
		if p.peekToken().Kind == tokens.TOKBigArrow {
			p.readToken()

			retValue, err := p.handleExpressions(tokens.PrecedenceLowest)
			if err != nil {
				return nil, err
			}

			blockNode.Value = retValue

			_, err = p.expect(tokens.TOKSemiColon, "expected ';' after block return statement")
			if err != nil {
				return nil, err
			}

			break
		}

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
				if err == nil {
					_, err = p.expect(tokens.TOKSemiColon, "expected ';' after expression statement")
				}
			}

		case tokens.TOKAt:
			statem, err = p.handleExpressions(tokens.PrecedenceLowest)
			if err == nil {
				_, err = p.expect(tokens.TOKSemiColon, "expected ';' after compiler action statement")
			}

		case tokens.TOKIf:
			statem, err = p.handleIfStatement()
		case tokens.TOKWhile:
			statem, err = p.handleWhileStatement()
		case tokens.TOKFor:
			statem, err = p.handleForStatement()

		default:
			statem, err = p.handleExpressions(tokens.PrecedenceLowest)
			if err == nil {
				_, err = p.expect(tokens.TOKSemiColon, "expected ';' after statement")
			}
		}

		if err != nil {
			return nil, err
		}

		if statem != nil {
			blockNode.Statements = append(blockNode.Statements, statem)
		}
	}

	_, err = p.expect(tokens.TOKRBraces, "expected '}' after function body")
	if err != nil {
		return nil, err
	}

	return blockNode, nil
}

// +-------------+
// | ControlFlow |
// +-------------+

func (p *parser) handleIfStatement() (ast.ASTNode, error) {
	p.readToken()

	ifNode := &ast.IfStatementNode{
		Conditions: []ast.IfCondition{},
		ElseBody:   nil,
	}
	_, err := p.expect(tokens.TOKLParen, "expected '(' before if condition")
	if err != nil {
		return nil, err
	}

	cond, err := p.handleExpressions(tokens.PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after if statement")
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	ifNode.Conditions = append(ifNode.Conditions, ast.IfCondition{
		Condition: cond,
		Body:      bodyBlock,
	})

	for p.peekToken().Kind == tokens.TOKElse {
		p.readToken()

		if p.peekToken().Kind == tokens.TOKIf {
			p.readToken()
			_, err := p.expect(tokens.TOKLParen, "expected '(' before if condition")
			if err != nil {
				return nil, err
			}

			nextCond, err := p.handleExpressions(tokens.PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			_, err = p.expect(tokens.TOKRParen, "expected ')' after if statement")
			if err != nil {
				return nil, err
			}

			nextBody, err := p.handleBlock()
			if err != nil {
				return nil, err
			}

			ifNode.Conditions = append(ifNode.Conditions, ast.IfCondition{
				Condition: nextCond,
				Body:      nextBody,
			})
		} else {
			elseBodyBlock, err := p.handleBlock()
			if err != nil {
				return nil, err
			}

			ifNode.ElseBody = elseBodyBlock
		}
	}

	return ifNode, nil
}

func (p *parser) handleWhileStatement() (ast.ASTNode, error) {
	p.readToken()

	_, err := p.expect(tokens.TOKLParen, "expected '(' before while condition")
	if err != nil {
		return nil, err
	}

	cond, err := p.handleExpressions(tokens.PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after while statement")
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatementNode{
		Condition: cond,
		Body:      bodyBlock,
	}, nil
}

func (p *parser) handleForStatement() (ast.ASTNode, error) {
	p.readToken()

	var init ast.ASTNode = nil
	var cond ast.ASTNode = nil
	var post ast.ASTNode = nil

	_, err := p.expect(tokens.TOKLParen, "expected '(' before for statement")
	if err != nil {
		return nil, err
	}

	if p.peekToken().Kind != tokens.TOKSemiColon {
		tok := p.peekToken()
		var err error

		if tok.Kind == tokens.TOKLet || tok.Kind == tokens.TOKConst {
			p.readToken()
			init, err = p.handleVariableDecl(false, tok.Kind == tokens.TOKConst)
		} else if tok.Kind == tokens.TOKId || p.peekNToken(1).Kind.IsAssign() {
			p.readToken()
			init, err = p.handleAssignement(tok)
		} else {
			init, err = p.handleExpressions(tokens.PrecedenceLowest)
		}
		if err != nil {
			return nil, err
		}
	}
	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after init of for statement")
	if err != nil {
		return nil, err
	}

	if p.peekToken().Kind != tokens.TOKSemiColon {
		cond, err = p.handleExpressions(tokens.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
	}
	_, err = p.expect(tokens.TOKSemiColon, "expected ';' after condition of for statement")
	if err != nil {
		return nil, err
	}

	if p.peekToken().Kind != tokens.TOKRParen {
		tok := p.peekToken()
		if p.peekNToken(1).Kind.IsAssign() {
			p.readToken()
			post, err = p.handleAssignement(tok)
		} else {
			post, err = p.handleExpressions(tokens.PrecedenceLowest)
		}
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(tokens.TOKRParen, "expected ')' after for statement")
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	return &ast.ForStatementNode{
		Init:      init,
		Condition: cond,
		Post:      post,
		Body:      bodyBlock,
	}, nil
}

// +---------------+
// | TypeExtension |
// +---------------+

func (p *parser) handleTypeExtension(targetType string) (ast.ASTNode, error) {
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
			p.readToken()
			identTok, errt := p.expect(tokens.TOKId, "expected compiler action name")
			if errt != nil {
				return nil, errt
			}

			switch identTok.Lexeme {
			case "overload":
				methodNode, err = p.handleCompileOverload()
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

			default:
				panic("not yet implemented the '@' (exluded the 'overload') symbol in the type extension")
			}

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

func (p *parser) handleCompileOverload() (ast.ASTNode, error) {
	_, err := p.expect(tokens.TOKAt, "expected '@' before action name in overload")
	if err != nil {
		return nil, err
	}

	actionTok, err := p.expect(tokens.TOKId, "expected action name after '@'")
	if err != nil {
		return nil, err
	}

	params, err := p.handleParameters()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(tokens.TOKArrow, "expected '->' after overload parameters")
	if err != nil {
		return nil, err
	}

	returnType, err := p.handleTypeSignature()
	if err != nil {
		return nil, err
	}

	if p.peekToken().Kind == tokens.TOKAt {
		p.readToken()
		invalidTok, err := p.expect(tokens.TOKId, "expected compiler directive")
		if err != nil || invalidTok.Lexeme != "invalid" {
			return nil, fmt.Errorf("expected 'invalid' directive after '@'")
		}

		return &ast.OverloadDeclNode{
			ActionName: actionTok.Lexeme,
			Parameters: params,
			ReturnType: returnType,
			Body:       nil,
			IsInvalid:  true,
		}, nil
	}

	body, err := p.handleBlock()
	if err != nil {
		return nil, err
	}

	return &ast.OverloadDeclNode{
		ActionName: actionTok.Lexeme,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
		IsInvalid:  false,
	}, nil
}
