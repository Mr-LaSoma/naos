package ast

type VariableDeclNode struct {
	IsPublic   bool
	IsConst    bool
	Name       string
	Type       ASTNode
	Expression ASTNode
}

func (n *VariableDeclNode) Node() {}
