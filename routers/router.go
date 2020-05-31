package routers

import (
	"github.com/PointStoneTeam/cugsync/api"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// 测试接口
	r.GET("/testapi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.GET("/getAllJob", api.GetAllJob)
	// 根据名称获取历史记录
	r.GET("/getHistory", api.GetHistory)
	// 获取文件列表
	r.GET("/getFolder", api.GetFolder)
	return r
}
