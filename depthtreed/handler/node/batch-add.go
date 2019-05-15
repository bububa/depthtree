package node

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

type BatchAddRequest = []AddRequest

func BatchAddHandler(c *gin.Context) {
	var req BatchAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	dbname := c.Param("db")
	if CheckWithCode(dbname == "", BADREQUEST_ERROR, "missing db", c) {
		return
	}
	treeBase := Service.TreeBase
	tree := treeBase.NewTree(dbname)
	for _, i := range req {
		tree.AddNode(i.Pid, i.Id)
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
