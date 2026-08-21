package ast

type ASTNode interface {
	Node()
	String() string
}

type ASTFile struct {
	PackageName string
	Imports     []ASTNode
	Decls       []ASTNode
}
