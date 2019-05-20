package router

import (
	"github.com/bububa/depthtree/depthtreed/handler/db"
	"github.com/gin-gonic/gin"
)

func dbRouter(r *gin.Engine) {
	dbGroup := r.Group("/db")
	dbGroup.POST("/truncate", db.TruncateHandler)
	dbGroup.GET("/roots/:db", db.RootsHandler)
	dbGroup.GET("/top-children/:db/:depth/:limit", db.TopChildrenHandler)
	dbGroup.GET("/list", db.ListHandler)
}
