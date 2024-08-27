package controllers

import (
	"loan-tracker-api/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogController struct {
    LogService domain.LogService
}

func NewLogController(logService domain.LogService) *LogController {
    return &LogController{LogService: logService}
}

func (lc *LogController) GetLogs(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	Type := c.Query("type")
	if Type == "" {
		Type = "all"
	}

	




    logs, err := lc.LogService.GetSystemLogs( Type, limit, page)
    if err.Message != "" {
        c.JSON(err.StatusCode, gin.H{
			"error": err.Message,
		})
        return
    }
    c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"message": "Logs retrieved successfully",
		"data": logs,

		"quantity": strconv.Itoa(len(logs)) + "/" + limit,
		"current_page": page,

	})
}