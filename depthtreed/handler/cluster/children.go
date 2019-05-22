package cluster

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	"github.com/bububa/depthtree"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
	"strconv"
)

func ChildrenHandler(c *gin.Context) {
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.Use(dbname)
	if CheckWithCode(tree == nil, NOTFOUND_ERROR, "not found db", c) {
		return
	}
	depth64, _ := strconv.ParseInt(c.Param("depth"), 10, 64)
	if CheckWithCode(depth64 == 0, BADREQUEST_ERROR, "missing depth", c) {
		return
	}
	k64, _ := strconv.ParseInt(c.Param("k"), 10, 64)
	if CheckWithCode(k64 == 0, BADREQUEST_ERROR, "missing k", c) {
		return
	}
	limit64, _ := strconv.ParseUint(c.Query("limit"), 10, 64)
	var (
		limit = int(limit64)
		depth = int(depth64)
		k     = int(k64)
	)
	clusters := tree.ChildrenCountInDepthCluster(depth, k)
	for _, c := range clusters {
		var nodes []*depthtree.Node
		for _, node := range c.Nodes {
			if len(nodes) >= limit {
				break
			}
			n := node.Copy(nil)
			n.ChildrenCount = int32(n.ChildrenCountInDepth(depth))
			nodes = append(nodes, n)
		}
		c.Nodes = nodes
		if len(c.Roots) > 100 {
			c.Roots = nil
		}
	}
	c.JSON(http.StatusOK, clusters)
}
