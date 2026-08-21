package ast

type VariableDeclNode struct {
	IsPublic   bool
	IsConst    bool
	Name       string
	Type       ASTNode
	Expression ASTNode
}

func (n *VariableDeclNode) Node() {}

// +-----------+
// | Functions |
// +-----------+

type Parameter struct {
	IsConst bool
	Name    string
	Type    ASTNode
}

type FuncDeclNode struct {
	IsPublic   bool
	Name       string
	Parameters []Parameter
	ReturnType ASTNode
	Body       []ASTNode
}

func (n *FuncDeclNode) Node() {}

type AssignementNode struct {
	Left     string
	Operator string
	Right    ASTNode
}

func (n *AssignementNode) Node() {}

type TypeExtensionNode struct {
	TargetType string
	Methods    []ASTNode
}

func (n *TypeExtensionNode) Node() {}
