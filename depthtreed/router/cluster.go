package router

import (
	"github.com/bububa/depthtree/depthtreed/handler/cluster"
	"github.com/gin-gonic/gin"
)

func clusterRouter(r *gin.Engine) {
	clusterGroup := r.Group("/cluster")
	clusterGroup.GET("/depth/:db/:k", cluster.DepthHandler)
	clusterGroup.GET("/children/:db/:depth/:k", cluster.ChildrenHandler)
}
