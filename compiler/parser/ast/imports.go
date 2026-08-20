package ast

import "fmt"

type ImportGlobalNode struct {
	Path string
}

func (n *ImportGlobalNode) Node() {}

func (n *ImportGlobalNode) String() string {
	return fmt.Sprintf("[global import] -> %s", n.Path)
}

type ImportAliasNode struct {
	Alias string
	Path  string
}

func (n *ImportAliasNode) Node() {}

func (n *ImportAliasNode) String() string {
	return fmt.Sprintf("[alias import '%s'] -> %s", n.Alias, n.Path)
}
