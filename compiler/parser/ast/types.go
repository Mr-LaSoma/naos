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
