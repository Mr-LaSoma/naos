package ast

import (
	"fmt"
	"strings"
)

// --- Expressions ---

func (n *PrefixExpressionNode) String() string {
	return fmt.Sprintf("(%s%s)", n.Operator, n.Right.String())
}

func (n *CallExpressionNode) String() string {
	args := make([]string, len(n.Arguments))
	for i, a := range n.Arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("%s(%s)", n.Callee.String(), strings.Join(args, ", "))
}

func (n *BinaryExpressionNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Operator, n.Right.String())
}

// --- Imports ---

func (n *ImportGlobalNode) String() string {
	return fmt.Sprintf("import %q", n.Path)
}

func (n *ImportAliasNode) String() string {
	return fmt.Sprintf("import %s %q", n.Alias, n.Path)
}

// --- Literals ---

func (n *IntLiteralNode) String() string {
	return n.Value
}

func (n *FloatLiteralNode) String() string {
	return n.Value
}

func (n *StringLiteralNode) String() string {
	return fmt.Sprintf("%q", n.Value)
}

func (n *CharLiteralNode) String() string {
	return fmt.Sprintf("'%s'", n.Value)
}

func (n *IdentifierRefNode) String() string {
	return n.Name
}

// --- Declarations ---

func (n *VariableDeclNode) String() string {
	var sb strings.Builder

	if n.IsPublic {
		sb.WriteString("pub ")
	}
	if n.IsConst {
		sb.WriteString("const ")
	} else {
		sb.WriteString("let ")
	}

	sb.WriteString(n.Name)

	if n.Type != nil {
		sb.WriteString(": ")
		sb.WriteString(n.Type.String())
	}

	if n.Expression != nil {
		sb.WriteString(" = ")
		sb.WriteString(n.Expression.String())
	}

	return sb.String()
}

func (n *TypeAliasNode) String() string {
	prefix := ""
	if n.IsPublic {
		prefix = "pub "
	}
	return fmt.Sprintf("%stype %s = %s", prefix, n.Name, n.SourceType.String())
}

func (n *NewTypeNode) String() string {
	prefix := ""
	if n.IsPublic {
		prefix = "pub "
	}
	return fmt.Sprintf("%stype %s %s", prefix, n.Name, n.BaseType.String())
}

// --- Types ---

func (n *TypeReferanceNode) String() string {
	return n.Name
}

func (n *PointerTypeNode) String() string {
	return fmt.Sprintf("*%s", n.ElementType.String())
}

func (n *MultiPointerTypeNode) String() string {
	return fmt.Sprintf("[*]%s", n.ElementType.String())
}

// --- Struct ---

func (f *StructField) String() string {
	prefix := ""
	if f.IsPublic {
		prefix = "pub "
	}
	return fmt.Sprintf("%s%s: %s", prefix, f.Name, f.Type.String())
}

func (n *StructNode) String() string {
	fields := make([]string, len(n.Fields))
	for i, f := range n.Fields {
		fields[i] = f.String()
	}
	return fmt.Sprintf("struct { %s }", strings.Join(fields, ", "))
}
