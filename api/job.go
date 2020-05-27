package api

import (
	"github.com/PointStoneTeam/cugsync/manager"
	"github.com/PointStoneTeam/cugsync/pkg/app"
	"github.com/PointStoneTeam/cugsync/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllJob(c *gin.Context) {
	jobList, err := manager.GetAllJobs()
	if err != nil {
		app.Response(c, http.StatusOK, e.JOB_GET_FAILED, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, jobList)
	}
}
