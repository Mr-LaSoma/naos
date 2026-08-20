package ast

import "fmt"

// +------------+
// | Base Types |
// +------------+

type TypeAliasNode struct {
	IsPublic   bool
	Name       string
	SourceType ASTNode
}

func (n *TypeAliasNode) Node() {}
func (n *TypeAliasNode) String() string {
	s := "private"
	if n.IsPublic {
		s = "public"
	}
	return fmt.Sprintf("[%s alias type %s] of {%v}", s, n.Name, n.SourceType)
}

type NewTypeNode struct {
	IsPublic bool
	Name     string
	BaseType ASTNode
}

func (n *NewTypeNode) Node() {}
func (n *NewTypeNode) String() string {
	s := "private"
	if n.IsPublic {
		s = "public"
	}
	return fmt.Sprintf("[%s new type %s] of {%v}", s, n.Name, n.BaseType)
}

type TypeReferanceNode struct {
	Name string
}

func (n *TypeReferanceNode) Node() {}
func (n *TypeReferanceNode) String() string {
	return "[type " + n.Name + "]"
}

// +----------+
// | Pointers |
// +----------+

type PointerTypeNode struct { // *type
	ElementType ASTNode
}

func (n *PointerTypeNode) Node() {}
func (n *PointerTypeNode) String() string {
	return fmt.Sprintf("[pointer] of {%v}", n.ElementType)
}

type MultiPointerTypeNode struct { // [*]type
	ElementType ASTNode
}

func (n *MultiPointerTypeNode) Node() {}
func (n *MultiPointerTypeNode) String() string {
	return fmt.Sprintf("[multi pointer] of {%v}", n.ElementType)
}

// +---------+
// | Structs |
// +---------+

type StructField struct {
	IsPublic bool
	Name     string
	Type     ASTNode
}

func (f StructField) String() string {
	s := "private"
	if f.IsPublic {
		s = "public"
	}
	return fmt.Sprintf("[%s field %s] of {%v}", s, f.Name, f.Type)
}

type StructNode struct {
	Fields []StructField
}

func (n *StructNode) Node() {}
func (n *StructNode) String() string {
	s := "[struct]{\n"
	for _, field := range n.Fields {
		s += fmt.Sprintf("\t%v;\n", field)
	}
	s += "}"
	return s
}
