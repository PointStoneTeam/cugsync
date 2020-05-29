package api

import (
	"github.com/PointStoneTeam/cugsync/folder"
	"github.com/PointStoneTeam/cugsync/pkg/app"
	"github.com/PointStoneTeam/cugsync/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFolder(c *gin.Context) {
	path := c.Query("path")

	list, err := folder.GetFolder(path)
	if err != nil {
		app.Response(c, http.StatusOK, e.FOLDER_GET_FAILED, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, list)
	}
}
