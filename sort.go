package depthtree

type DepthNodeSlice []*Node

func (c DepthNodeSlice) Len() int {
	return len(c)
}

func (c DepthNodeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c DepthNodeSlice) Less(i, j int) bool {
	if c[i].MaxDepth == c[j].MaxDepth {
		return c[i].Id < c[j].Id
	}
	return c[i].MaxDepth < c[j].MaxDepth
}

type ChildrenCountNodeSlice struct {
	nodes []*Node
	depth int
}

func NewChildrenCountNodeSlice(nodes []*Node, depth int) ChildrenCountNodeSlice {
	return ChildrenCountNodeSlice{
		nodes: nodes,
		depth: depth,
	}
}

func (c ChildrenCountNodeSlice) Nodes() []*Node {
	return c.nodes
}

func (c ChildrenCountNodeSlice) Len() int {
	return len(c.nodes)
}

func (c ChildrenCountNodeSlice) Swap(i, j int) {
	c.nodes[i], c.nodes[j] = c.nodes[j], c.nodes[i]
}

func (c ChildrenCountNodeSlice) Less(i, j int) bool {
	iCount := c.nodes[i].ChildrenCountInDepth(c.depth)
	jCount := c.nodes[j].ChildrenCountInDepth(c.depth)
	if iCount == jCount {
		return c.nodes[i].Id < c.nodes[j].Id
	}
	return iCount < jCount
}
