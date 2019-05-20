package depthtree

import (
	//"fmt"
	"github.com/bububa/depthtree/clusters"
	"math"
	"sort"
)

type DepthClusterPoint struct {
	Node *Node
}

func (c DepthClusterPoint) Coordinates() clusters.Coordinates {
	return clusters.Coordinates([]float64{float64(c.Node.MaxDepth), 0})
}

// Distance returns the euclidean distance between two coordinates
func (c DepthClusterPoint) Distance(p2 clusters.Coordinates) float64 {
	var r float64
	p1 := c.Coordinates()
	for i, v := range p1 {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

func (c DepthClusterPoint) Data() interface{} {
	return c.Node
}

type ChildrenCountClusterPoint struct {
	Node  *Node
	Depth int
}

func (c ChildrenCountClusterPoint) Coordinates() clusters.Coordinates {
	var childrenCount = float64(c.Node.ChildrenCountInDepth(c.Depth))
	return clusters.Coordinates([]float64{childrenCount, 0})
}

// Distance returns the euclidean distance between two coordinates
func (c ChildrenCountClusterPoint) Distance(p2 clusters.Coordinates) float64 {
	var r float64
	p1 := c.Coordinates()
	for i, v := range p1 {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

func (c ChildrenCountClusterPoint) Data() interface{} {
	return c.Node
}

type Cluster struct {
	Center clusters.Coordinates `json:"center"`
	Range  []int                `json:"range"`
	Count  int                  `json:"count"`
	Nodes  []*Node              `json:"nodes,omitempty"`
	Roots  []*Node              `json:"roots,omitempty"`
}

func NewCluster(cluster clusters.Cluster, depth int) *Cluster {
	this := &Cluster{
		Center: cluster.Center,
		Count:  len(cluster.Observations),
	}
	var (
		min     float64 = -1
		max     float64
		nodes   []*Node
		nodeMap = make(map[int64]*Node, this.Count)
	)
	for _, o := range cluster.Observations {
		p := o.Coordinates()
		if min == -1 || min > p[0] {
			min = p[0]
		}
		if max < p[0] {
			max = p[0]
		}
		node := o.Data().(*Node)
		nodes = append(nodes, node)
		nodeMap[node.Id] = node
	}
	if depth == -2 {
		sliceSorter := DepthNodeSlice(nodes)
		sort.Sort(sort.Reverse(sliceSorter))
		this.Nodes = sliceSorter
	} else {
		sliceSorter := NewChildrenCountNodeSlice(nodes, depth)
		sort.Sort(sort.Reverse(sliceSorter))
		this.Nodes = sliceSorter.Nodes()
	}
	tree := NewTree()
	for _, node := range this.Nodes {
		var foundParent bool
		for _, p := range node.Parents() {
			if _, found := nodeMap[p.Id]; found {
				tree.AddNode(p.Id, node.Id)
				foundParent = true
			}
		}
		if !foundParent {
			tree.AddRoot(node.Id)
		}
	}
	roots := tree.RootNodes()
	for _, node := range roots {
		if depth == -2 {
			n := node.Copy(nil)
			if ori, found := nodeMap[n.Id]; found {
				n.MaxDepth = ori.MaxDepth
				n.MinDepth = ori.MinDepth
			}
			this.Roots = append(this.Roots, n)
		} else {
			if ori, found := nodeMap[node.Id]; found {
				this.Roots = append(this.Roots, ori)
			}
		}

	}
	if depth != -2 {
		sliceSorter := NewChildrenCountNodeSlice(this.Roots, depth)
		sort.Sort(sort.Reverse(sliceSorter))
		rootNodes := sliceSorter.Nodes()
		this.Roots = []*Node{}
		for _, node := range rootNodes {
			n := node.Copy(nil)
			n.ChildrenCount = node.ChildrenCountInDepth(depth)
			this.Roots = append(this.Roots, n)
		}
	}
	this.Range = []int{int(min), int(max)}
	return this
}
