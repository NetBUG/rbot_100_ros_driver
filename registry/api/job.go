package api

import (
	"encoding/json"
	//"fmt"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	//"github.com/franela/goreq"
	"github.com/gin-gonic/gin"
	"net/http"
)

// apps/:app_id/jobs
func (api APIBase) GetJobList(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var js []model.Job
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

// apps/:app_id/jobs
func (api APIBase) CreateJob(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var job model.Job

	id := c.Param("app_id")
	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	/*var njob model.Job
	if e := api.DB.Model(&f).Association("Jobs").Append([]model.Job{njob}).Error; e != nil {
		errors.Append(APIError{Message: "Could not save the code."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}


	c.JSON(http.StatusOK, job)
	return*/

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Action      string `json:"action" valid:"optional"`
		Type        string `json:"type" valid:"optional"`
		Parameters  string `json:"parameters" valid:"optional"`
		Input       string `json:"input" valid:"optional"`
		Output      string `json:"output" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()

	job.Name = req.Name
	job.Description = req.Description
	job.Input = req.Input
	job.Output = req.Output
	job.Action = req.Action
	job.Type = req.Type
	job.Parameters = req.Parameters
	//job.Appid = f.Appid;

	// validation against the request schema
	/*if _, err := valid.ValidateStruct(req); err != nil {
		errors.Append(APIError{Message: err.Error()})
		c.JSON(http.StatusNotFound, errors)
		return
	}*/

	//TBD: Check this logic
	/*var js []model.Job
	if e := api.DB.Model(&f).Related(&js).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}*/
	//End of TBD

	//TBD: Atomic transaction
	if e := api.DB.Create(&job).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	/*
		if f.Jobs != "" {
			f.Jobs+=","+job.ID;
		} else {
			f.Jobs=job.ID;
		}
		f.Jobs = append(f.Jobs, job.Jobid);
		if e := api.DB.Model(&f).Updates(f).Error; e != nil {
			errors.Append(APIError{Message: "Database error"})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}*/
	//End of TBD
	if e := api.DB.Model(&f).Association("Jobs").Append([]model.Job{job}).Error; e != nil {
		errors.Append(APIError{Message: "Could not save the code."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	c.JSON(http.StatusOK, job)
}

// /jobs/:id
func (api APIBase) GetJob(c *gin.Context) {
	var errors APIErrors
	var job model.Job

	id := c.Param("job_id")
	if e := api.DB.Find(&job, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.JSON(http.StatusOK, job)
}

// /jobs/:id
func (api APIBase) UpdateJob(c *gin.Context) {
	var errors APIErrors
	var job model.Job
	var f model.App

	id := c.Param("job_id")
	if e := api.DB.Find(&job, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Could not found the job"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Find(&f).Where("id = ?", job.AppID).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Action      string `json:"action" valid:"optional"`
		Type        string `json:"type" valid:"optional"`
		Parameters  string `json:"parameters" valid:"optional"`
		Input       string `json:"input" valid:"optional"`
		Output      string `json:"output" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()
	job.Name = req.Name
	job.Description = req.Description
	job.Input = req.Input
	job.Output = req.Output
	job.Action = req.Action
	job.Type = req.Type
	job.Parameters = req.Parameters

	if e := api.DB.Model(&job).Updates(&job).Error; e != nil {
		errors.Append(APIError{Message: "Could not update the job."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, job)
}

// /jobs/:id
func (api APIBase) DeleteJob(c *gin.Context) {
	var errors APIErrors
	var job model.Job

	id := c.Param("job_id")
	if e := api.DB.Find(&job, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	//TBD: Delete reference to job in app
	/*
		if e := api.DB.Find(&f).Where("id = ?", job.AppID).Error; e != nil {
			errors.Append(APIError{Message: "Database error"})
			c.JSON(http.StatusNotFound, errors)
			return
		}
		jobs = f.Jobs.split(",");
		for j in jobs) {
			if j = job.ID) {
				continue;
			}
			njobs+=","+j;
		}
		f.Jobs = njobs;
		if e := api.DB.Model(&f).Updates(&f).Error; e != nil {
			errors.Append(APIError{Message: "Could not update the job."})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}*/

	//Delete job
	if e := api.DB.Delete(&job).Error; e != nil {
		errors.Append(APIError{Message: "Data deletion error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "The job was deleted successfully.",
	})
}
