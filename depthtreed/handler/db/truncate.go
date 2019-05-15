package db

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

type TruncateRequest struct {
	Name string `json:"name" binding:"required"`
}

func TruncateHandler(c *gin.Context) {
	var req TruncateRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	treeBase := Service.TreeBase
	treeBase.Truncate(req.Name)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
