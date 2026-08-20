package nodes

type GlobalImportNode struct {
	Path string
}

func (n *GlobalImportNode) Node() {}

type AliasImportNode struct {
	Alias string
	Path  string
}

func (n *AliasImportNode) Node() {}
