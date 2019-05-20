package depthtree

import (
	"fmt"
	"github.com/xlab/treeprint"
	"strconv"
	"sync"
)

type Node struct {
	Id                   int64              `json:"i"`
	parents              []*Node            `json:"-"`
	Children             []*Node            `json:"c,omitempty"`
	MaxDepth             int                `json:"mxd,omitempty"`
	MinDepth             int                `json:"mnd,omitempty"`
	childrenCountInDepth map[int]int        `json:"-"`
	ChildrenCount        int                `json:"n,omitempty"`
	parentMap            map[int64]struct{} `json:"-"`
	childrenMap          map[int64]struct{} `json:"-"`
	sync.RWMutex
}

func NewNode(id int64) *Node {
	return &Node{
		Id:                   id,
		parentMap:            make(map[int64]struct{}),
		childrenMap:          make(map[int64]struct{}),
		childrenCountInDepth: make(map[int]int),
	}
}

func (this *Node) String() string {
	return strconv.FormatInt(this.Id, 10)
}

func (this *Node) HasParent() bool {
	this.RLock()
	defer this.RUnlock()
	return len(this.parents) > 0
}

func (this *Node) Parents() []*Node {
	this.RLock()
	defer this.RUnlock()
	var nodes []*Node
	for _, node := range this.parents {
		nodes = append(nodes, node)
	}
	return nodes
}

func (this *Node) HasChildren() bool {
	this.RLock()
	defer this.RUnlock()
	return len(this.Children) > 0
}

func (this *Node) ChildrenIds() []int64 {
	this.RLock()
	defer this.RUnlock()
	var ids []int64
	for _, node := range this.Children {
		ids = append(ids, node.Id)
	}
	return ids
}

func (this *Node) ChildrenCountInDepth(depth int) int {
	this.RLock()
	defer this.RUnlock()
	if c, found := this.childrenCountInDepth[depth]; found {
		return c
	}
	return 0
}

func (this *Node) RemoveChild(nodeId int64) {
	this.RLock()
	var (
		nodes []*Node
		found bool
	)
	for _, node := range this.Children {
		if node.Id == nodeId {
			found = true
			continue
		}
		nodes = append(nodes, node)
	}
	this.RUnlock()
	this.Lock()
	if found {
		delete(this.childrenMap, nodeId)
	}
	this.Children = nodes
	this.Unlock()
}

func (this *Node) RemoveFromChildren() {
	this.RLock()
	for _, node := range this.Children {
		node.RemoveParent(this.Id)
	}
	this.RUnlock()
}

func (this *Node) RemoveParent(nodeId int64) {
	this.RLock()
	var (
		nodes []*Node
		found bool
	)
	for _, node := range this.parents {
		if node.Id == nodeId {
			found = true
			continue
		}
		nodes = append(nodes, node)
	}
	this.RUnlock()
	this.Lock()
	if found {
		delete(this.parentMap, nodeId)
	}
	this.parents = nodes
	this.Unlock()
}

func (this *Node) AddParent(node *Node) {
	this.RLock()
	if _, found := this.parentMap[node.Id]; found {
		this.RUnlock()
		return
	}
	this.RUnlock()
	this.Lock()
	this.parents = append(this.parents, node)
	this.parentMap[node.Id] = struct{}{}
	this.Unlock()
}

func (this *Node) AddChild(node *Node) {
	this.RLock()
	if _, found := this.childrenMap[node.Id]; found {
		this.RUnlock()
		return
	}
	this.RUnlock()
	this.Lock()
	this.Children = append(this.Children, node)
	this.childrenMap[node.Id] = struct{}{}
	this.Unlock()
	node.AddParent(this)
}

func (this *Node) AddChildren(nodes []*Node) {
	this.Lock()
	for _, node := range nodes {
		if _, found := this.childrenMap[node.Id]; found {
			continue
		}
		this.Children = append(this.Children, node)
		this.childrenMap[node.Id] = struct{}{}
	}
	this.Unlock()
	for _, node := range nodes {
		node.AddParent(this)
	}
}

func (this *Node) Copy(children []*Node) *Node {
	node := NewNode(this.Id)
	this.RLock()
	node.MaxDepth = this.MaxDepth
	node.MinDepth = this.MinDepth
	this.RUnlock()
	node.AddChildren(children)
	return node
}

func (this *Node) Depth() (int, int) {
	this.RLock()
	children := this.Children
	this.RUnlock()
	if len(children) == 0 {
		this.Lock()
		this.MaxDepth = 0
		this.MinDepth = 0
		this.Unlock()
		return 0, 0
	}
	var (
		maxDepth int
		minDepth int
	)
	for _, node := range children {
		minD, maxD := node.Depth()
		if maxDepth < maxD {
			maxDepth = maxD
		}
		if minDepth > minD {
			minDepth = minD
		}
	}
	this.Lock()
	this.MaxDepth = maxDepth + 1
	this.MinDepth = minDepth + 1
	this.Unlock()
	return minDepth + 1, maxDepth + 1
}

func (this *Node) GetChildren(depth int) ([]*Node, int) {
	this.RLock()
	children := this.Children
	totalChildren := len(children)
	this.RUnlock()
	if depth < 0 {
		depth = -1
	}
	if depth == 0 || totalChildren == 0 {
		this.Lock()
		this.childrenCountInDepth[depth] = totalChildren
		this.Unlock()
		return nil, totalChildren
	}
	var nodes []*Node
	for _, child := range children {
		childNodes, childrenCount := child.GetChildren(depth - 1)
		totalChildren += childrenCount
		node := child.Copy(childNodes)
		node.ChildrenCount = childrenCount
		nodes = append(nodes, node)
	}
	this.Lock()
	this.childrenCountInDepth[depth] = totalChildren
	this.ChildrenCount = totalChildren
	this.Unlock()
	return nodes, totalChildren
}

func (this *Node) CountChildren(depth int) int {
	this.RLock()
	children := this.Children
	totalChildren := len(children)
	this.RUnlock()
	if depth < 0 {
		depth = -1
	}
	if depth == 0 || totalChildren == 0 {
		this.Lock()
		this.childrenCountInDepth[depth] = totalChildren
		this.Unlock()
		return totalChildren
	}
	for _, child := range children {
		childCount := child.CountChildren(depth - 1)
		totalChildren += childCount
	}
	this.Lock()
	this.childrenCountInDepth[depth] = totalChildren
	this.Unlock()
	return totalChildren
}

func (this *Node) PrintTree(depth int) string {
	printTree := treeprint.New()
	this.printTreeIter(printTree, depth)
	return printTree.String()
}

func (this *Node) printTreeIter(printTree treeprint.Tree, depth int) {
	if depth < 0 {
		depth = -1
	}
	this.RLock()
	key := this.String()
	var childrenCount int
	if c, found := this.childrenCountInDepth[depth]; found {
		childrenCount = c
	}
	meta := fmt.Sprintf("maxD: %d, minD: %d, childCount: %d", this.MaxDepth, this.MinDepth, childrenCount)
	children := this.Children
	this.RUnlock()
	if !this.HasChildren() {
		printTree.AddMetaNode(meta, key)
		return
	}
	branch := printTree.AddMetaBranch(meta, key)
	for _, child := range children {
		child.printTreeIter(branch, depth-1)
	}
}
