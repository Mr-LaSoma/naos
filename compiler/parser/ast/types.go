package ast

// +------------+
// | Base Types |
// +------------+

type TypeAliasNode struct {
	IsPublic   bool
	Name       string
	SourceType ASTNode
}

func (n *TypeAliasNode) Node() {}

type NewTypeNode struct {
	IsPublic bool
	Name     string
	BaseType ASTNode
}

func (n *NewTypeNode) Node() {}

type TypeReferanceNode struct {
	Name string
}

func (n *TypeReferanceNode) Node() {}

// +----------+
// | Pointers |
// +----------+

type PointerTypeNode struct { // *type
	ElementType ASTNode
}

func (n *PointerTypeNode) Node() {}

type MultiPointerTypeNode struct { // [*]type
	ElementType ASTNode
}

func (n *MultiPointerTypeNode) Node() {}

// +---------+
// | Structs |
// +---------+

type StructField struct {
	IsPublic bool
	Name     string
	Type     ASTNode
}

type StructNode struct {
	Fields []StructField
}

func (n *StructNode) Node() {}

type FuncSignatureNode struct {
	IsPublic   bool
	Name       string
	Parameter  []Parameter
	ReturnType ASTNode
}

func (n *FuncSignatureNode) Node() {}

type InterfaceLiteralNode struct {
	Methods []*FuncSignatureNode
}

func (n *InterfaceLiteralNode) Node() {}

type EnumVariant struct {
	Name   string
	Fields []StructField
}

type EnumLiteralNode struct {
	Variants []EnumVariant
}

func (n *EnumLiteralNode) Node() {}

type GenericTypeNode struct {
	Name       string
	Constraint ASTNode
}

func (n *GenericTypeNode) Node() {}

type GenericInstantiationNode struct {
	Left         ASTNode
	TypeArgument []ASTNode
}

func (n *GenericInstantiationNode) Node() {}
