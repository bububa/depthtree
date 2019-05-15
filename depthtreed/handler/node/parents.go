package node

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/bububa/depthtree"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
	"strconv"
)

func ParentsHandler(c *gin.Context) {
	nodeId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if CheckWithCode(nodeId == 0, BADREQUEST_ERROR, "missing node id", c) {
		return
	}
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.Use(dbname)
	if CheckWithCode(tree == nil, NOTFOUND_ERROR, "not found db", c) {
		return
	}
	node := tree.Find(nodeId)
	if CheckWithCode(node == nil, NOTFOUND_ERROR, "not found node", c) {
		return
	}
	nodes := node.Parents()
	var parents []*depthtree.Node
	for _, n := range nodes {
		parents = append(parents, n.Copy(nil))
	}
	c.JSON(http.StatusOK, parents)
}
