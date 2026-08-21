package ast

type ImportGlobalNode struct {
	Path string
}

func (n *ImportGlobalNode) Node() {}

type ImportAliasNode struct {
	Alias string
	Path  string
}

func (n *ImportAliasNode) Node() {}
