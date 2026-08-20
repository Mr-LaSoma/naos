package parser

import "naoslang/tokenizer/tokens"

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

func (p *parser) Parse() (*ASTFile, error) {
	err := p.dropUselessAndErrors()
	if err != nil {
		return nil, err
	}

	fileAst := &ASTFile{
		Imports: []ASTNode{},
		Decls:   []ASTNode{},
	}
	err = p.handlePackageDecl(fileAst)
	if err != nil {
		return nil, err
	}

	err = p.handleImports(fileAst)
	if err != nil {
		return nil, err
	}

	return fileAst, nil
}
