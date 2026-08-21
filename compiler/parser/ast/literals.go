package ast

// +---------+
// | Numbers |
// +---------+

type IntLiteralNode struct {
	Value string
}

func (n *IntLiteralNode) Node() {}

type FloatLiteralNode struct {
	Value string
}

func (n *FloatLiteralNode) Node() {}

// +---------+
// | Strings |
// +---------+

type StringLiteralNode struct {
	Value string
}

func (n *StringLiteralNode) Node() {}

type CharLiteralNode struct {
	Value string
}

func (n *CharLiteralNode) Node() {}

// +------+
// | Misc |
// +------+

type IdentifierRefNode struct {
	Name string
}

func (n *IdentifierRefNode) Node() {}
