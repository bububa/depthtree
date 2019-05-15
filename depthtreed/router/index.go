package router

import (
	"github.com/bububa/depthtree/depthtreed/common"
	"github.com/bububa/depthtree/depthtreed/router/static"
	"github.com/danielkov/gin-helmet"
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(templatePath string, config common.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(helmet.NoSniff(), helmet.DNSPrefetchControl(), helmet.SetHSTS(true), helmet.IENoOpen(), helmet.XSSFilter())
	xssMdlwr := &xss.XssMw{
		FieldsToSkip: []string{"password", "start_date", "end_date", "token"},
		BmPolicy:     "UGCPolicy",
	}
	r.Use(xssMdlwr.RemoveXss())
	r.Use(static.Serve("/", static.LocalFile(config.StaticPath, 0, true)))
	r.LoadHTMLGlob(templatePath)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "depthtreed"})
		return
	})
	nodeRouter(r)
	dbRouter(r)
	clusterRouter(r)
	return r
}
