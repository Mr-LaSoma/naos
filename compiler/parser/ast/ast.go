package ast

type ASTNode interface {
	Node()
}

type ASTFile struct {
	PackageName string
	Imports     []ASTNode
	Decls       []ASTNode
}
