package api

import (
	"github.com/PointStoneTeam/cugsync/folder"
	"github.com/PointStoneTeam/cugsync/pkg/app"
	"github.com/PointStoneTeam/cugsync/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @title 获取文件夹列表
// @version v0.0.0
// @Summary getFolder
// @Description 获取文件夹列表
// @Description ```code 10002 获取文件夹列表失败```
// @Param name query string true "centos"
// @Accept  json
// @Produce  json
// @Success 200 {object} app.Res{data=object{sum=integer,list=[]folder.File}}
// @Router /getFolder [get]
func GetFolder(c *gin.Context) {
	path := c.Query("path")

	list, err := folder.GetFolder(path)
	if err != nil {
		app.Response(c, http.StatusOK, e.FOLDER_GET_FAILED, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, list)
	}
}
