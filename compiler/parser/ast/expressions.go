package ast

type PrefixExpressionNode struct {
	Operator string
	Right    ASTNode
}

func (n *PrefixExpressionNode) Node() {}

type CallExpressionNode struct {
	Callee    ASTNode
	Arguments []ASTNode
}

func (n *CallExpressionNode) Node() {}

type BinaryExpressionNode struct {
	Left     ASTNode
	Operator string
	Right    ASTNode
}

func (n *BinaryExpressionNode) Node() {}
