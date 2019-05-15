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

func DepthHandler(c *gin.Context) {
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.Use(dbname)
	if CheckWithCode(tree == nil, NOTFOUND_ERROR, "not found db", c) {
		return
	}
	k64, _ := strconv.ParseInt(c.Param("k"), 10, 64)
	if CheckWithCode(k64 == 0, BADREQUEST_ERROR, "missing k", c) {
		return
	}
	limit64, _ := Uint64Value(c.Query("limit"), 10)
	var (
		k     = int(k64)
		limit = int(limit64)
	)
	clusters := tree.DepthCluster(k)
	for _, c := range clusters {
		var nodes []*depthtree.Node
		for _, node := range c.Nodes {
			if len(nodes) >= limit {
				break
			}
			nodes = append(nodes, node.Copy(nil))
		}
		c.Nodes = nodes
	}
	c.JSON(http.StatusOK, clusters)
}
