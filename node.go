package depthtree

import (
	"fmt"
	"github.com/xlab/treeprint"
	"strconv"
	"sync"
	"sync/atomic"
)

type NodeMeta interface {
	String() string
}

type Node struct {
	Id              int64                   `json:"i"`
	parents         []*Node                 `json:"-"`
	Children        []*Node                 `json:"c,omitempty"`
	MaxDepth        int32                   `json:"mxd,omitempty"`
	MinDepth        int32                   `json:"mnd,omitempty"`
	childrenInDepth *NodeChildrenDepthCount `json:"-"`
	ChildrenCount   int32                   `json:"n,omitempty"`
	parentMap       map[int64]struct{}      `json:"-"`
	childrenMap     map[int64]struct{}      `json:"-"`
	Meta            NodeMeta                `json:"meta,omitempty"`
	sync.RWMutex
}

type NodeChildrenDepthCount struct {
	mp map[int]int
	sync.RWMutex
}

func NewNodeChildrenDepthCount() *NodeChildrenDepthCount {
	return &NodeChildrenDepthCount{
		mp: make(map[int]int),
	}
}

func (this *NodeChildrenDepthCount) Set(depth int, count int) {
	this.Lock()
	this.mp[depth] = count
	this.Unlock()
}

func (this *NodeChildrenDepthCount) Get(depth int) int {
	this.RLock()
	this.RUnlock()
	if count, found := this.mp[depth]; found {
		return count
	}
	return 0
}

func NewNode(id int64) *Node {
	return &Node{
		Id:              id,
		parentMap:       make(map[int64]struct{}),
		childrenMap:     make(map[int64]struct{}),
		childrenInDepth: NewNodeChildrenDepthCount(),
	}
}

func (this *Node) String() string {
	return strconv.FormatInt(this.Id, 10)
}

func (this *Node) HasParent() bool {
	this.RLock()
	hasParent := len(this.parents) > 0
	this.RUnlock()
	return hasParent
}

func (this *Node) Parents() []*Node {
	this.RLock()
	nodes := make([]*Node, len(this.parents))
	copy(nodes, this.parents)
	this.RUnlock()
	return nodes
}

func (this *Node) HasChildren() bool {
	this.RLock()
	hasChildren := len(this.Children) > 0
	this.RUnlock()
	return hasChildren
}

func (this *Node) ChildrenIds() []int64 {
	var ids []int64
	this.RLock()
	for _, node := range this.Children {
		ids = append(ids, node.Id)
	}
	this.RUnlock()
	return ids
}

func (this *Node) ChildrenCountInDepth(depth int) int {
	return this.childrenInDepth.Get(depth)
}

func (this *Node) RemoveChild(nodeId int64) {
	var (
		nodes []*Node
		found bool
	)
	this.RLock()
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	for _, node := range children {
		if node.Id == nodeId {
			found = true
			continue
		}
		nodes = append(nodes, node)
	}
	this.Lock()
	if found {
		delete(this.childrenMap, nodeId)
	}
	this.Children = nodes
	this.Unlock()
}

func (this *Node) RemoveFromChildren() {
	this.RLock()
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	for _, node := range children {
		node.RemoveParent(this.Id)
	}
}

func (this *Node) RemoveParent(nodeId int64) {
	var (
		nodes []*Node
		found bool
	)
	this.RLock()
	parents := make([]*Node, len(this.parents))
	copy(parents, this.parents)
	this.RUnlock()
	for _, node := range parents {
		if node.Id == nodeId {
			found = true
			continue
		}
		nodes = append(nodes, node)
	}
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
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	if len(children) == 0 {
		atomic.StoreInt32(&this.MaxDepth, 0)
		atomic.StoreInt32(&this.MinDepth, 0)
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
	atomic.StoreInt32(&this.MaxDepth, int32(maxDepth+1))
	atomic.StoreInt32(&this.MinDepth, int32(minDepth+1))
	return minDepth + 1, maxDepth + 1
}

func (this *Node) GetChildren(depth int) ([]*Node, int) {
	this.RLock()
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	totalChildren := len(children)
	if depth < 0 {
		depth = -1
	}
	if depth == 0 || totalChildren == 0 {
		this.childrenInDepth.Set(depth, totalChildren)
		return nil, totalChildren
	}
	var nodes []*Node
	for _, child := range children {
		childNodes, childrenCount := child.GetChildren(depth - 1)
		totalChildren += childrenCount
		node := child.Copy(childNodes)
		node.ChildrenCount = int32(childrenCount)
		nodes = append(nodes, node)
	}
	this.childrenInDepth.Set(depth, totalChildren)
	atomic.StoreInt32(&this.ChildrenCount, int32(totalChildren))
	return nodes, totalChildren
}

func (this *Node) CountChildren(depth int) int {
	this.RLock()
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	totalChildren := len(children)
	if depth < 0 {
		depth = -1
	}
	if depth == 0 || totalChildren == 0 {
		this.childrenInDepth.Set(depth, totalChildren)
		return totalChildren
	}
	for _, child := range children {
		childCount := child.CountChildren(depth - 1)
		totalChildren += childCount
	}
	this.childrenInDepth.Set(depth, totalChildren)
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
	key := this.String()
	childrenCount := this.childrenInDepth.Get(depth)
	meta := fmt.Sprintf("maxD: %d, minD: %d, childCount: %d", this.MaxDepth, this.MinDepth, childrenCount)
	if !this.HasChildren() {
		printTree.AddMetaNode(meta, key)
		return
	}
	branch := printTree.AddMetaBranch(meta, key)
	this.RLock()
	children := make([]*Node, len(this.Children))
	copy(children, this.Children)
	this.RUnlock()
	for _, child := range children {
		child.printTreeIter(branch, depth-1)
	}
}
