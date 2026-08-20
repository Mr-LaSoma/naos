package parser

import (
	"naoslang/parser/nodes"
	"naoslang/tokenizer/tokens"
)

// +----------+
// | Packages |
// +----------+

func (p *parser) handlePackageDecl(fileAst *ASTFile) error {
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

func (p *parser) handleImports(fileAst *ASTFile) error {
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

func (p *parser) handleGlobalImport() (ASTNode, error) {
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
	return &nodes.GlobalImportNode{Path: cleanPath}, nil
}

func (p *parser) handleAliasImport(alias string) (ASTNode, error) {
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
	return &nodes.AliasImportNode{Alias: alias, Path: cleanPath}, nil
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
