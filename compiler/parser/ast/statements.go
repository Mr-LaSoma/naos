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

type FuncAttribute struct {
	Name string
	Args map[string]string
}

type FuncDeclNode struct {
	IsPublic   bool
	Name       string
	Parameters []Parameter
	ReturnType ASTNode
	Body       ASTNode
	Attributes []FuncAttribute
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

type OverloadDeclNode struct {
	ActionName string
	Parameters []Parameter
	ReturnType ASTNode
	Body       ASTNode
	IsInvalid  bool
}

func (n *OverloadDeclNode) Node() {}

// +------------+
// | Conditions |
// +------------+

type IfCondition struct {
	Condition ASTNode
	Body      ASTNode
}

type IfStatementNode struct {
	Conditions []IfCondition
	ElseBody   ASTNode
}

func (n *IfStatementNode) Node() {}

type WhileStatementNode struct {
	Condition ASTNode
	Body      ASTNode
}

func (n *WhileStatementNode) Node() {}

type ForStatementNode struct {
	Init      ASTNode
	Condition ASTNode
	Post      ASTNode
	Body      ASTNode
}

func (n *ForStatementNode) Node() {}

type ReturnStatement struct {
	Expression ASTNode
}

func (n *ReturnStatement) Node() {}

type BreakStatementNode struct{}

func (n *BreakStatementNode) Node() {}

type ContinueStatementNode struct{}

func (n *ContinueStatementNode) Node() {}

type DeferStatementNode struct {
	Body ASTNode
}

func (n *DeferStatementNode) Node() {}

type MatchCase struct {
	Patterns []ASTNode
	Body     ASTNode
}

type MatchStatementNode struct {
	Expression ASTNode
	Cases      []MatchCase
}

func (n *MatchStatementNode) Node() {}
