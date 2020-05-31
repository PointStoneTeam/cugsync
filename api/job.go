package api

import (
	"github.com/PointStoneTeam/cugsync/manager"
	"github.com/PointStoneTeam/cugsync/pkg/app"
	"github.com/PointStoneTeam/cugsync/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @title 获取所有任务
// @version v0.0.0
// @Summary getAllJob
// @Description 获取所有任务接口
// @Description ```code 10000 任务列表获取失败```
// @Accept  json
// @Produce  json
// @Success 200 {object} app.Res{data=object{sum=integer,list=[]manager.UnCreatedJob}}
// @Router /getAllJob [get]
func GetAllJob(c *gin.Context) {
	data := make(map[string]interface{})

	jobList, err := manager.GetAllJobs()
	if err != nil {
		app.Response(c, http.StatusOK, e.JOB_GET_FAILED, nil)
	} else {
		data["sum"] = len(jobList)
		data["list"] = jobList
		app.Response(c, http.StatusOK, e.SUCCESS, data)
	}
}

// @title 获取历史记录
// @version v0.0.0
// @Summary getHistory
// @Description 获取历史记录
// @Description ```code 10001 历史列表获取失败```
// @Param name query string true "centos"
// @Accept  json
// @Produce  json
// @Success 200 {object} app.Res{data=object{sum=integer,list=[]manager.History}}
// @Router /getHistory [get]
func GetHistory(c *gin.Context) {
	name := c.Query("name")

	data := make(map[string]interface{})
	if historyList, err := manager.GetHistory(name); err != nil {
		app.Response(c, http.StatusOK, e.HISTORY_GET_FAILED, nil)
	} else {
		data["sum"] = len(historyList)
		data["list"] = historyList
		app.Response(c, http.StatusOK, e.SUCCESS, data)
	}
}
