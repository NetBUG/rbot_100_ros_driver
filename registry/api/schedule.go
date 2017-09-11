package api

import (
	//"fmt"
	"encoding/json"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	//"github.com/franela/goreq"
	"github.com/gin-gonic/gin"
	"net/http"
)

// apps/:app_id/schedule
func (api APIBase) GetScheduleList(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var js []model.Schedule
	id := c.Param("app_id")

	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Model(&f).Related(&js).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, js)
}

// apps/:app_id/schedule
func (api APIBase) CreateSchedule(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var schedule model.Schedule

	id := c.Param("app_id")
	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Schedule    string `json:"action" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()

	schedule.Name = req.Name
	schedule.Description = req.Description
	schedule.Sched = req.Schedule
	schedule.AppID = f.ID

	// validation against the request schema
	/*if _, err := valid.ValidateStruct(req); err != nil {
		errors.Append(APIError{Message: err.Error()})
		c.JSON(http.StatusNotFound, errors)
		return
	}*/

	//TBD: Check this logic
	/*var js []model.Schedule
	if e := api.DB.Model(&f).Related(&js).Error; e != nil {
		errors.Append(APIError{Message: "Database error."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}*/
	//End of TBD

	//TBD: Atomic transaction
	if e := api.DB.Create(&schedule).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	/*f.Schedule=schedule.Scheduleid;
	if e := api.DB.Model(&f).Updates(f).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}*/
	if e := api.DB.Model(&f).Association("Schedules").Append([]model.Schedule{schedule}).Error; e != nil {
		errors.Append(APIError{Message: "Could not save the code."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	//End of TBD
	c.JSON(http.StatusOK, schedule)
}

// /schedule/:id
func (api APIBase) GetSchedule(c *gin.Context) {
	var errors APIErrors
	var schedule model.Schedule

	id := c.Param("schedule_id")
	if e := api.DB.Find(&schedule, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// /schedule/:id
func (api APIBase) UpdateSchedule(c *gin.Context) {
	var errors APIErrors
	var schedule model.Schedule
	var f model.App

	id := c.Param("schedule_id")
	if e := api.DB.Find(&schedule, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Find(&f).Where("id = ?", schedule.AppID).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Schedule    string `json:"schedule" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()
	schedule.Name = req.Name
	schedule.Description = req.Description
	schedule.Sched = req.Schedule

	if e := api.DB.Model(&schedule).Updates(&schedule).Error; e != nil {
		errors.Append(APIError{Message: "Data update error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// /schedule/:id
func (api APIBase) DeleteSchedule(c *gin.Context) {
	var errors APIErrors
	var schedule model.Schedule

	id := c.Param("schedule_id")
	if e := api.DB.Find(&schedule, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	//TBD: Delete reference to schedule in app
	/*
		if e := api.DB.Find(&f).Where("id = ?", schedule.AppID).Error; e != nil {
			errors.Append(APIError{Message: "Database error"})
			c.JSON(http.StatusNotFound, errors)
			return
		}
		schedule = f.Schedule.split(",");
		for j in schedule) {
			if j = schedule.ID) {
				continue;
			}
			nschedule+=","+j;
		}
		f.Schedule = nschedule;
		if e := api.DB.Model(&f).Updates(&f).Error; e != nil {
			errors.Append(APIError{Message: "Data update error"})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}*/

	//Delete schedule
	if e := api.DB.Delete(&schedule).Error; e != nil {
		errors.Append(APIError{Message: "Data deletion error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "The schedule was deleted successfully.",
	})
}
