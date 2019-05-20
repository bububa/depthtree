package db

import (
	//"github.com/davecgh/go-spew/spew"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
	"strconv"
)

func TopChildrenHandler(c *gin.Context) {
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
	limit64, _ := strconv.ParseUint(c.Param("limit"), 10, 64)
	var (
		limit = int(limit64)
		depth = int(depth64)
	)
	nodes := tree.TopChildrenCounts(depth, limit)
	c.JSON(http.StatusOK, nodes)
}
