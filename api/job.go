package api

import (
	"github.com/PointStoneTeam/cugsync/manager"
	"github.com/PointStoneTeam/cugsync/pkg/app"
	"github.com/PointStoneTeam/cugsync/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

// get all job
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

// get history by name
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
