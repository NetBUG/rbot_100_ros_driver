package api

import (
	"encoding/json"
	"fmt"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (api APIBase) GetAppList(c *gin.Context) {
	var errors APIErrors
	var fl []model.App

	if e := api.DB.Find(&fl).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	if len(fl) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, fl)
}
func (api APIBase) CreateApp(c *gin.Context) {
	var errors APIErrors
	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		User        string `json:"user" valid:"optional"`
	}
	var a model.App

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	defer c.Request.Body.Close()

	// validation against the request schema
	if _, err := valid.ValidateStruct(req); err != nil {
		errors.Append(APIError{Message: err.Error()})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	a.Name = req.Name
	a.Description = req.Description
	//TODO: Remove: Temporary fix until FE sets the user name
	//a.User = req.User
	a.User = "ciscozeus"
	if e := api.DB.Create(&a).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (api APIBase) GetApp(c *gin.Context) {
	var errors APIErrors
	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
	}

	var a model.App

	id := c.Param("app_id")
	if e := api.DB.Find(&a, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.JSON(http.StatusOK, a)
}

func (api APIBase) UpdateApp(c *gin.Context) {
	var errors APIErrors
	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		User        string `json:"user" valid:"optional"`
	}

	var a model.App

	id := c.Param("app_id")
	if e := api.DB.Find(&a, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		fmt.Printf("Database error \n")
		c.JSON(http.StatusNotFound, errors)
		return
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		fmt.Printf("Error json parse \n")
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()

	// validation against the request schema
	if _, err := valid.ValidateStruct(req); err != nil {
		fmt.Printf("Error schema validation : \n")
		errors.Append(APIError{Message: err.Error()})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	a.Name = req.Name
	a.Description = req.Description
	a.User = req.User

	if e := api.DB.Model(&a).Updates(a).Error; e != nil {
		fmt.Printf("Error model update : \n")
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (api APIBase) DeleteApp(c *gin.Context) {
	var errors APIErrors
	var a model.App
	id := c.Param("app_id")

	if e := api.DB.Find(&a, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Delete(&a).Error; e != nil {
		// This might not be happend because if DB connection is losted, the previous SQL would be failed.
		errors.Append(APIError{Message: "App deletion error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "The app was deleted successfully.",
	})
}

func (api APIBase) CloneApp(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var nf model.App

	id := c.Param("app_id")
	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error: Could not find app"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	nf.Name = f.Name
	nf.Description = f.Description
	//TBD: Update User
	//Create another app with the new parameters.
	if e := api.DB.Create(&nf).Error; e != nil {
		errors.Append(APIError{Message: "Database error: Could not clone app"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	//Find jobs
	var jl []model.Job
	if e := api.DB.Find(&jl, "app_id=?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error: Could not find job"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	//Data input for one job, may be the output of another job, so
	//we need to avoid cloning the data entry multiple times in the loops below
	//Use a hash table to prevent it from being cloned multiple time.
	var datacreated map[string]int
	datacreated = make(map[string]int)
	if len(jl) != 0 {
		//Clone jobs

		for _, job := range jl {
			/*
				if e := api.DB.Find(&j, "id = ?", job).Error; e != nil {
					errors.Append(APIError{Message: "Database error"})
					c.JSON(http.StatusNotFound, errors)
					return
				}
			*/
			var inputdata []string
			var outputdata []string
			var j model.Job
			j.Name = job.Name
			j.Description = job.Description
			j.Action = job.Action
			j.Parameters = job.Parameters
			inputdata = strings.Split(job.Input, ",")
			outputdata = strings.Split(job.Output, ",")

			//Create data for job
			for _, dataid := range inputdata {
				var d model.Data
				if e := api.DB.Find(&d, "id = ?", dataid).Error; e != nil {
					errors.Append(APIError{Message: "Database error: Could not find data"})
					c.JSON(http.StatusNotFound, errors)
					return
				}

				//Check if data has already been created. Data input from one job, may be the output of another job,
				//so the data may already have been created.
				if datacreated[dataid] == 0 {
					//Create new data
					if e := api.DB.Create(&d).Error; e != nil {
						errors.Append(APIError{Message: "Database error: Could not clone data"})
						c.JSON(http.StatusInternalServerError, errors)
						return
					}
					if e := api.DB.Model(&nf).Association("Data").Append([]model.Data{d}).Error; e != nil {
						errors.Append(APIError{Message: "Could not save the data."})
						c.JSON(http.StatusInternalServerError, errors)
						return
					}
					datacreated[dataid] = 1
				}

				//Update new job with new data
				if j.Input != "" {
					j.Input += "," + d.ID
				} else {
					j.Input = d.ID
				}

			}
			for _, dataid := range outputdata {
				var d model.Data
				if e := api.DB.Find(&d, "id = ?", dataid).Error; e != nil {
					errors.Append(APIError{Message: "Database error: Could not find data"})
					c.JSON(http.StatusNotFound, errors)
					return
				}

				//Create new data
				//Check if data has already been created. Data input from one job, may be the output of another job,
				//so the data may already have been created.
				if datacreated[dataid] == 0 {
					if e := api.DB.Create(&d).Error; e != nil {
						errors.Append(APIError{Message: "Database error: Could not clone data"})
						c.JSON(http.StatusInternalServerError, errors)
						return
					}
					if e := api.DB.Model(&nf).Association("Data").Append([]model.Data{d}).Error; e != nil {
						errors.Append(APIError{Message: "Could not save the data."})
						c.JSON(http.StatusInternalServerError, errors)
						return
					}
					datacreated[dataid] = 1
				}
				//Update new job with new data
				if j.Output != "" {
					j.Output += "," + d.ID
				} else {
					j.Output = d.ID
				}

			}
			//Clone job
			if e := api.DB.Create(&j).Error; e != nil {
				errors.Append(APIError{Message: "Database error: Could not clone job"})
				c.JSON(http.StatusInternalServerError, errors)
				return
			}
			if e := api.DB.Model(&nf).Association("Jobs").Append([]model.Job{j}).Error; e != nil {
				errors.Append(APIError{Message: "Could not save the job."})
				c.JSON(http.StatusInternalServerError, errors)
				return
			}
		}
	}

	//Find schedule
	var sl []model.Schedule
	if e := api.DB.Find(&sl, "app_id=?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error: Could not find job"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	//Clone schedule
	if len(sl) != 0 {
		for _, schedule := range sl {
			var s model.Schedule
			s.Name = schedule.Name
			s.Description = schedule.Description
			s.Sched = schedule.Sched
			if e := api.DB.Create(&s).Error; e != nil {
				errors.Append(APIError{Message: "Database error: Could not clone schedule"})
				c.JSON(http.StatusInternalServerError, errors)
				return
			}
			if e := api.DB.Model(&nf).Association("Schedules").Append([]model.Schedule{s}).Error; e != nil {
				errors.Append(APIError{Message: "Could not save the schedule."})
				c.JSON(http.StatusInternalServerError, errors)
				return
			}
		}
	}

	//Provision resources needed by app here
	c.JSON(http.StatusOK, f)
}
