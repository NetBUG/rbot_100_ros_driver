package api

import (
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	"github.com/gin-gonic/gin"
	"net/http"
	//	"strconv"
)

// APIRoot
// Establish connection to DB and make Client
func APIRoot() (*gin.Engine, error) {
	db, err := model.OpenDB()
	if err != nil {
		return nil, err
	}

	base := APIBase{
		DB: db,
	}

	g := gin.Default()
	g.Use(gin.Recovery())

	// GET /
	g.GET("", base.GetRoot)

	// /apps
	appGroup := g.Group("apps")
	appGroup.POST("", base.CreateApp)
	appGroup.GET("/:app_id", base.GetApp)
	appGroup.POST("/:app_id", base.UpdateApp)
	appGroup.DELETE("/:app_id", base.DeleteApp)
	appGroup.GET("", base.GetAppList)
	// jobs
	appGroup.POST("/:app_id/job", base.CreateJob)
	appGroup.GET("/:app_id/jobs", base.GetJobList)
	jobGroup := g.Group("job")
	jobGroup.GET("/:job_id", base.GetJob)
	jobGroup.POST("/:job_id", base.UpdateJob)
	jobGroup.DELETE("/:job_id", base.DeleteJob)
	// data
	appGroup.POST("/:app_id/data", base.CreateData)
	appGroup.GET("/:app_id/data", base.GetDataList)
	dataGroup := g.Group("data")
	dataGroup.GET("/:data_id", base.GetData)
	dataGroup.POST("/:data_id", base.UpdateData)
	dataGroup.DELETE("/:data_id", base.DeleteData)
	// schedule
	appGroup.POST("/:app_id/schedule", base.CreateSchedule)
	appGroup.GET("/:app_id/schedule", base.GetScheduleList)
	scheduleGroup := g.Group("schedule")
	scheduleGroup.GET("/:schedule_id", base.GetSchedule)
	scheduleGroup.POST("/:schedule_id", base.UpdateSchedule)
	scheduleGroup.DELETE("/:schedule_id", base.DeleteSchedule)

	//Cloning a app
	cloneGroup := g.Group("clone")
	cloneGroup.POST("/:app_id", base.CloneApp)
	return g, nil
}

func (api APIBase) GetRoot(c *gin.Context) {
	c.JSON(http.StatusOK, DefaultResponse{
		Message: "Endpoint works successfully.",
	})
}
