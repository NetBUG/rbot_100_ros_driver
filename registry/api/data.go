package api

import (
	"encoding/json"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	//"github.com/franela/goreq"
	"github.com/gin-gonic/gin"
	"net/http"
)

// apps/:app_id/data
func (api APIBase) GetDataList(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var ds []model.Data
	id := c.Param("app_id")

	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "App not found"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Model(&f).Related(&ds).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	c.JSON(http.StatusOK, ds)
}

// apps/:app_id/data
func (api APIBase) CreateData(c *gin.Context) {
	var errors APIErrors
	var f model.App
	var data model.Data

	id := c.Param("app_id")
	if e := api.DB.Find(&f, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "App not found"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Action      string `json:"action" valid:"optional"`
		Type        string `json:"type" valid:"optional"`
		Store       string `json:"store" valid:"optional"`
		Parameters  string `json:"parameters" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON request parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()

	data.Name = req.Name
	data.Description = req.Description
	data.Type = req.Type
	data.Store = req.Store
	data.Parameters = req.Parameters

	// validation against the request schema
	/*if _, err := valid.ValidateStruct(req); err != nil {
		errors.Append(APIError{Message: err.Error()})
		c.JSON(http.StatusNotFound, errors)
		return
	}*/

	//TBD: Check this logic
	/*	var js []model.Data
		if e := api.DB.Model(&f).Related(&js).Error; e != nil {
			errors.Append(APIError{Message: "Data not found"})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
	*/ //End of TBD

	//TBD: Atomic transaction
	if e := api.DB.Create(&data).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	/*
		if f.Data != "" {
			f.Data+=","+data.ID;
		} else {
			f.Data=data.ID;
		}
		f.Data = append(f.Data, data.Dataid);
		if e := api.DB.Model(&f).Updates(f).Error; e != nil {
			errors.Append(APIError{Message: "Something wrong with database."})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}
	*/
	if e := api.DB.Model(&f).Association("Data").Append([]model.Data{data}).Error; e != nil {
		errors.Append(APIError{Message: "Could not save the code."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}
	//End of TBD
	c.JSON(http.StatusOK, data)
}

// /data/:id
func (api APIBase) GetData(c *gin.Context) {
	var errors APIErrors
	var data model.Data

	id := c.Param("data_id")
	if e := api.DB.Find(&data, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "App not found"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.JSON(http.StatusOK, data)
}

// /data/:id
func (api APIBase) UpdateData(c *gin.Context) {
	var errors APIErrors
	var data model.Data
	var f model.App

	id := c.Param("data_id")
	if e := api.DB.Find(&data, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "Database error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	if e := api.DB.Find(&f).Where("id = ?", data.AppID).Error; e != nil {
		errors.Append(APIError{Message: "App not found"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	type Request struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"optional"`
		Action      string `json:"action" valid:"optional"`
		Type        string `json:"type" valid:"optional"`
		Store       string `json:"store" valid:"optional"`
		Parameters  string `json:"parameters" valid:"optional"`
	}

	var req Request
	decoder := json.NewDecoder(c.Request.Body)

	if e := decoder.Decode(&req); e != nil {
		errors.Append(APIError{Message: "JSON parse error"})
		c.JSON(http.StatusNotFound, errors)
		return
	}
	c.Request.Body.Close()
	data.Name = req.Name
	data.Description = req.Description
	data.Type = req.Type
	data.Store = req.Store
	data.Parameters = req.Parameters

	if e := api.DB.Model(&data).Updates(&data).Error; e != nil {
		errors.Append(APIError{Message: "Data update error."})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, data)
}

// /data/:id
func (api APIBase) DeleteData(c *gin.Context) {
	var errors APIErrors
	var data model.Data

	id := c.Param("data_id")
	if e := api.DB.Find(&data, "id = ?", id).Error; e != nil {
		errors.Append(APIError{Message: "App not found"})
		c.JSON(http.StatusNotFound, errors)
		return
	}

	//TBD: Delete reference to data in job, app
	/*
		if e := api.DB.Find(&f).Where("id = ?", data.AppID).Error; e != nil {
			errors.Append(APIError{Message: "App not found"})
			c.JSON(http.StatusNotFound, errors)
			return
		}
		datas = f.Data.split(",");
		for j in datas) {
			if j = data.ID) {
				continue;
			}
			ndatas+=","+j;
		}
		f.Data = ndata;
		if e := api.DB.Model(&f).Updates(&f).Error; e != nil {
			errors.Append(APIError{Message: "Data update error."})
			c.JSON(http.StatusInternalServerError, errors)
			return
		}*/

	//Delete data
	if e := api.DB.Delete(&data).Error; e != nil {
		errors.Append(APIError{Message: "Data deletion error"})
		c.JSON(http.StatusInternalServerError, errors)
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Message: "The data was deleted successfully.",
	})
}
