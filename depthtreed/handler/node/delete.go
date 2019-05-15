package node

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

type DeleteRequest struct {
	Id int64 `json:"id" binding:"required"`
}

func DeleteHandler(c *gin.Context) {
	var req DeleteRequest
	if CheckErr(c.Bind(&req), c) {
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
	tree.RemoveNode(req.Id)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
