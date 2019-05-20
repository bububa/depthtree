package router

import (
	"github.com/bububa/depthtree/depthtreed/handler/node"
	"github.com/gin-gonic/gin"
)

func nodeRouter(r *gin.Engine) {
	nodeGroup := r.Group("/node")
	nodeGroup.GET("/children/:db/:id/:depth", node.ChildrenHandler)
	nodeGroup.GET("/info/:db/:id", node.InfoHandler)
	nodeGroup.GET("/parents/:db/:id", node.ParentsHandler)
	nodeGroup.POST("/add/:db", node.AddHandler)
	nodeGroup.POST("/batch-add/:db", node.BatchAddHandler)
	nodeGroup.POST("/delete/:db", node.DeleteHandler)
}
