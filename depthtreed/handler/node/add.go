package node

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

type AddRequest struct {
	Id  int64 `json:"id" binding:"required"`
	Pid int64 `json:"pid" binding:"required"`
}

func AddHandler(c *gin.Context) {
	var req AddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.NewTree(dbname)
	tree.AddNode(req.Pid, req.Id)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
