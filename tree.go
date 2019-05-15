package depthtree

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/bububa/depthtree/clusters"
	"github.com/bububa/depthtree/kmeans"
	"github.com/mkideal/log"
	"io"
	"os"
	"sync"
)

func Int64ToBytes(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes(), err
}

func BytesToInt64(b []byte) (int64, error) {
	buf := bytes.NewBuffer(b)
	var i int64
	err := binary.Read(buf, binary.BigEndian, &i)
	return i, err
}

type Tree struct {
	nodeMap map[int64]*Node `json:"-"`
	sync.RWMutex
}

func NewTree() *Tree {
	return &Tree{
		nodeMap: make(map[int64]*Node),
	}
}

func (this *Tree) AddRoot(id int64) {
	this.RLock()
	if _, found := this.nodeMap[id]; found {
		this.RUnlock()
		return
	}
	this.RUnlock()
	node := NewNode(id)
	this.Lock()
	this.nodeMap[id] = node
	this.Unlock()
}

func (this *Tree) AddNode(from int64, to int64) {
	var (
		fromNode *Node
		toNode   *Node
		newFrom  bool
		newTo    bool
	)
	this.RLock()
	if node, found := this.nodeMap[from]; found {
		fromNode = node
	} else {
		fromNode = NewNode(from)
		newFrom = true
	}
	if node, found := this.nodeMap[to]; found {
		toNode = node
	} else {
		toNode = NewNode(to)
		newTo = true
	}
	this.RUnlock()
	if newFrom || newTo {
		this.Lock()
		if newFrom {
			this.nodeMap[from] = fromNode
		}
		if newTo {
			this.nodeMap[to] = toNode
		}
		this.Unlock()
	}
	fromNode.AddChild(toNode)
}

func (this *Tree) RemoveNode(nodeId int64) bool {
	node := this.Find(nodeId)
	if node == nil {
		return false
	}
	parents := node.Parents()
	for _, node := range parents {
		node.RemoveChild(nodeId)
	}
	node.RemoveFromChildren()
	this.Lock()
	delete(this.nodeMap, nodeId)
	this.Unlock()
	return true
}

func (this *Tree) Find(nodeId int64) *Node {
	this.RLock()
	defer this.RUnlock()
	if node, found := this.nodeMap[nodeId]; found {
		return node
	}
	return nil
}

func (this *Tree) RootNodes() []*Node {
	this.RLock()
	var roots []*Node
	for _, node := range this.nodeMap {
		if node.HasParent() {
			continue
		}
		roots = append(roots, node)
	}
	this.RUnlock()
	return roots
}

func (this *Tree) DepthCluster(k int) []*Cluster {
	var points clusters.Observations
	roots := this.RootNodes()
	for _, node := range roots {
		node.Depth()
	}
	this.RLock()
	for _, node := range this.nodeMap {
		points = append(points, DepthClusterPoint{Node: node})
	}
	this.RUnlock()
	km := kmeans.New()
	clusters, err := km.Partition(points, k)
	if err != nil {
		panic(err)
	}
	var ret []*Cluster
	for _, c := range clusters {
		cluster := NewCluster(c, -2)
		ret = append(ret, cluster)
	}
	return ret
}

func (this *Tree) ChildrenCountInDepthCluster(depth int, k int) []*Cluster {
	var points clusters.Observations
	roots := this.RootNodes()
	for _, node := range roots {
		node.Depth()
		if depth < 0 {
			node.CountChildren(depth)
		}
	}
	wg := &sync.WaitGroup{}
	this.RLock()
	for _, node := range this.nodeMap {
		if depth >= 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, n *Node) {
				n.CountChildren(depth)
				wg.Done()
			}(wg, node)
		}
		points = append(points, ChildrenCountClusterPoint{Node: node, Depth: depth})
	}
	this.RUnlock()
	wg.Wait()
	km := kmeans.New()
	clusters, err := km.Partition(points, k)
	if err != nil {
		panic(err)
	}
	var ret []*Cluster
	for _, c := range clusters {
		cluster := NewCluster(c, depth)
		ret = append(ret, cluster)
	}
	return ret
}

func (this *Tree) Flush(name string) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	this.RLock()
	w := bufio.NewWriter(f)
	for _, node := range this.nodeMap {
		pid, err := Int64ToBytes(node.Id)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		ids := node.ChildrenIds()
		for _, i := range ids {
			id, err := Int64ToBytes(i)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			w.Write(pid)
			w.Write(id)
		}
	}
	this.RUnlock()
	w.Flush()
	return nil
}

func (this *Tree) Reload(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	line := make([]byte, 16)
	for {
		_, err := bfRd.Read(line)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err.Error())
			break
		}
		if len(line) != 16 {
			log.Warn("invalid line: %d", len(line))
			continue
		}
		pid, err := BytesToInt64(line[0:8])
		if err != nil {
			log.Error(err.Error())
			continue
		}
		nid, err := BytesToInt64(line[8:16])
		if err != nil {
			log.Error(err.Error())
			continue
		}
		this.AddNode(pid, nid)
	}
	return nil
}
