package db

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/bububa/depthtree"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

func RootsHandler(c *gin.Context) {
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.Use(dbname)
	if CheckWithCode(tree == nil, NOTFOUND_ERROR, "not found db", c) {
		return
	}
	nodes := tree.RootNodes()
	var roots []*depthtree.Node
	for _, n := range nodes {
		n.Depth()
		childCount := n.CountChildren(-1)
		node := n.Copy(nil)
		node.ChildrenCount = int32(childCount)
		roots = append(roots, node)
	}
	c.JSON(http.StatusOK, roots)
}
