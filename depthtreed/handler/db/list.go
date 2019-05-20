package db

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/bububa/depthtree/depthtreed/common"
	. "github.com/bububa/depthtree/depthtreed/handler"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"net/http"
)

func ListHandler(c *gin.Context) {
	treeBase := Service.TreeBase
	dbs := treeBase.List()
	c.JSON(http.StatusOK, dbs)
}
