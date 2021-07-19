package routers

import (
	"github.com/gin-gonic/gin"
	"save2table/controllers"
	_ "save2table/docs"
	"save2table/middleware"
)

func SetUp() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.Cors())
	//r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/v1/pics", controllers.PicsHandler)
	r.GET("/index", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "老弟"})
	})
	return r
}
