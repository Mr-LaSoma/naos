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

// --- Functions ---

func (p *Parameter) String() string {
	prefix := ""
	if p.IsConst {
		prefix = "const "
	}
	return fmt.Sprintf("%s%s: %s", prefix, p.Name, p.Type.String())
}

func (n *FuncDeclNode) String() string {
	var sb strings.Builder

	if n.IsPublic {
		sb.WriteString("pub ")
	}
	sb.WriteString("func ")
	sb.WriteString(n.Name)
	sb.WriteString("(")

	params := make([]string, len(n.Parameters))
	for i := range n.Parameters {
		params[i] = n.Parameters[i].String()
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(")")

	if n.ReturnType != nil {
		sb.WriteString(" ")
		sb.WriteString(n.ReturnType.String())
	}

	sb.WriteString(" {")
	for _, stmt := range n.Body {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(stmt.String()))
	}
	if len(n.Body) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}

func indentBody(s string) string {
	lines := strings.Split(s, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = "\t" + lines[i]
	}
	return strings.Join(lines, "\n")
}

func (n *AssignementNode) String() string {
	return fmt.Sprintf("%s %s %s", n.Left, n.Operator, n.Right.String())
}

func (n *TypeExtensionNode) String() string {
	var sb strings.Builder

	sb.WriteString(n.TargetType)
	sb.WriteString(" # {")
	for _, m := range n.Methods {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(m.String()))
	}
	if len(n.Methods) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}
