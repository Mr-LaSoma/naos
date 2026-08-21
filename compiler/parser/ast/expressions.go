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

type CompilerActionNode struct {
	Name      string
	Left      ASTNode
	Arguments []ASTNode
}

func (n *CompilerActionNode) Node() {}

type MemberAccessNode struct {
	Left   ASTNode
	Member string
}

func (n *MemberAccessNode) Node() {}

type BlockExpressionNode struct {
	Statements []ASTNode
	Value      ASTNode
}

func (n *BlockExpressionNode) Node() {}

type IndexAccessNode struct {
	Left  ASTNode
	Index ASTNode
}

func (n *IndexAccessNode) Node() {}
