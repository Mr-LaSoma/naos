package ast

import (
	"fmt"
	"sort"
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

func (a *FuncAttribute) String() string {
	if len(a.Args) == 0 {
		return fmt.Sprintf("@%s", a.Name)
	}

	keys := make([]string, 0, len(a.Args))
	for k := range a.Args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	args := make([]string, len(keys))
	for i, k := range keys {
		args[i] = fmt.Sprintf("%s: %s", k, a.Args[k])
	}

	return fmt.Sprintf("@%s(%s)", a.Name, strings.Join(args, ", "))
}

func (n *FuncDeclNode) String() string {
	var sb strings.Builder

	for i := range n.Attributes {
		sb.WriteString(n.Attributes[i].String())
		sb.WriteString("\n")
	}

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
	sb.WriteString(n.Body.String())
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

func (n *CompilerActionNode) String() string {
	var sb strings.Builder
	if n.Left != nil {
		sb.WriteString(n.Left.String())
		sb.WriteString(" ")
	}
	sb.WriteString("@")
	sb.WriteString(n.Name)
	sb.WriteString("(")
	args := make([]string, len(n.Arguments))
	for i := range n.Arguments {
		args[i] = n.Arguments[i].String()
	}
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")")

	return sb.String()
}

func (n *OverloadDeclNode) String() string {
	var sb strings.Builder

	sb.WriteString("@overload @")
	sb.WriteString(n.ActionName)
	sb.WriteString("(")

	params := make([]string, len(n.Parameters))
	for i := range n.Parameters {
		params[i] = n.Parameters[i].String()
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(")")

	if n.ReturnType != nil {
		sb.WriteString(" -> ")
		sb.WriteString(n.ReturnType.String())
	}

	if n.IsInvalid {
		sb.WriteString(" @invalid")
		return sb.String()
	}

	sb.WriteString(" {")
	sb.WriteString(n.Body.String())
	sb.WriteString("}")

	return sb.String()
}

func (n *MemberAccessNode) String() string {
	return fmt.Sprintf("%v.%s", n.Left, n.Member)
}

func (n *FuncSignatureNode) String() string {
	var sb strings.Builder

	if n.IsPublic {
		sb.WriteString("pub ")
	}
	sb.WriteString(n.Name)
	sb.WriteString("(")

	params := make([]string, len(n.Parameter))
	for i := range n.Parameter {
		params[i] = n.Parameter[i].String()
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(")")

	if n.ReturnType != nil {
		sb.WriteString(" -> ")
		sb.WriteString(n.ReturnType.String())
	}

	return sb.String()
}

func (n *InterfaceLiteralNode) String() string {
	var sb strings.Builder

	sb.WriteString("interface {")
	for i := range n.Methods {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(n.Methods[i].String()))
	}
	if len(n.Methods) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}

func (v *EnumVariant) String() string {
	if len(v.Fields) == 0 {
		return v.Name
	}

	fields := make([]string, len(v.Fields))
	for i := range v.Fields {
		fields[i] = v.Fields[i].String()
	}
	return fmt.Sprintf("%s(%s)", v.Name, strings.Join(fields, ", "))
}

func (n *EnumLiteralNode) String() string {
	var sb strings.Builder

	sb.WriteString("enum {")
	for i := range n.Variants {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(n.Variants[i].String()))
		sb.WriteString(",")
	}
	if len(n.Variants) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}

func (n *GenericTypeNode) String() string {
	if n.Constraint == nil {
		return fmt.Sprintf("$%s", n.Name)
	}
	return fmt.Sprintf("$%s %s", n.Name, n.Constraint.String())
}

func (n *GenericInstantiationNode) String() string {
	args := make([]string, len(n.TypeArgument))
	for i, a := range n.TypeArgument {
		args[i] = a.String()
	}
	return fmt.Sprintf("%s(%s)", n.Left.String(), strings.Join(args, ", "))
}

func (n *BlockExpressionNode) String() string {
	if len(n.Statements) == 0 && n.Value == nil {
		return "{}"
	}

	var sb strings.Builder

	sb.WriteString("{")
	for _, stmt := range n.Statements {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(stmt.String()))
	}
	if n.Value != nil {
		sb.WriteString("\n\t=> ")
		sb.WriteString(indentBody(n.Value.String()))
	}
	sb.WriteString("\n}")

	return sb.String()
}

// --- Control flow ---

func (c *IfCondition) String() string {
	return fmt.Sprintf("%s %s", c.Condition.String(), c.Body.String())
}

func (n *IfStatementNode) String() string {
	var sb strings.Builder

	for i, cond := range n.Conditions {
		if i == 0 {
			sb.WriteString("if ")
		} else {
			sb.WriteString(" else if ")
		}
		sb.WriteString(cond.String())
	}

	if n.ElseBody != nil {
		sb.WriteString(" else ")
		sb.WriteString(n.ElseBody.String())
	}

	return sb.String()
}

func (n *WhileStatementNode) String() string {
	return fmt.Sprintf("while %s %s", n.Condition.String(), n.Body.String())
}

func (n *ForStatementNode) String() string {
	init, cond, post := "", "", ""
	if n.Init != nil {
		init = n.Init.String()
	}
	if n.Condition != nil {
		cond = n.Condition.String()
	}
	if n.Post != nil {
		post = n.Post.String()
	}

	return fmt.Sprintf("for %s; %s; %s %s", init, cond, post, n.Body.String())
}

func (n *ReturnStatement) String() string {
	if n.Expression == nil {
		return "return"
	}
	return fmt.Sprintf("return %s", n.Expression.String())
}

func (n *BreakStatementNode) String() string {
	return "break"
}

func (n *ContinueStatementNode) String() string {
	return "continue"
}

func (n *DeferStatementNode) String() string {
	return fmt.Sprintf("defer %s", n.Body.String())
}

// --- Match ---

func (c *MatchCase) String() string {
	patterns := make([]string, len(c.Patterns))
	for i, p := range c.Patterns {
		patterns[i] = p.String()
	}
	return fmt.Sprintf("%s => %s", strings.Join(patterns, ", "), c.Body.String())
}

func (n *MatchStatementNode) String() string {
	var sb strings.Builder

	sb.WriteString("match ")
	sb.WriteString(n.Expression.String())
	sb.WriteString(" {")
	for i := range n.Cases {
		sb.WriteString("\n\t")
		sb.WriteString(indentBody(n.Cases[i].String()))
	}
	if len(n.Cases) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}
