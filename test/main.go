package main

import (
	"fmt"
	"github.com/bububa/depthtree"
	//"github.com/davecgh/go-spew/spew"
)

func main() {
	tests := [][]int64{
		{1, 2},
		{1, 3},
		{3, 4},
		{3, 5},
		{5, 6},
		{2, 7},
		{2, 15},
		{7, 8},
		{6, 9},
		{6, 10},
		{9, 11},
		{9, 12},
		{9, 13},
		{13, 14},
	}
	tree := depthtree.NewTree()
	for _, t := range tests {
		tree.AddNode(t[0], t[1])
	}
	roots := tree.RootNodes()
	for _, node := range roots {
		node.Depth()
		node.CountChildren(2)
		fmt.Println(node.PrintTree(2))
		// children := node.GetChildren(2)
		// spew.Dump(children)
	}
	// tree.DepthCluster(3)
	clusters := tree.ChildrenCountInDepthCluster(2, 3)
	for _, c := range clusters {
		fmt.Println("======")
		for _, n := range c.Nodes {
			fmt.Printf("n: %s, c: %d\n", n.String(), n.ChildrenCountInDepth(2))
		}
	}
}
