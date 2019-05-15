package router

import (
	"github.com/bububa/depthtree/depthtreed/handler/db"
	"github.com/gin-gonic/gin"
)

func dbRouter(r *gin.Engine) {
	dbGroup := r.Group("/db")
	dbGroup.POST("/truncate", db.TruncateHandler)
	dbGroup.GET("/roots/:db", db.RootsHandler)
}
