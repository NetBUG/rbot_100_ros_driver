package api

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type APIBase struct {
	DB *gorm.DB
	//Wiz *wiz.Client
}

type DefaultResponse struct {
	Message string
}

type APIErrors struct {
	Errors []APIError
}

func (es *APIErrors) Append(e APIError) []APIError {
	es.Errors = append(es.Errors, e)
	return es.Errors
}

func (es *APIErrors) Len() int {
	return len(es.Errors)
}

type APIError struct {
	Message string `json:"message"`
}

// Serve runs API server process
func Serve() error {
	r, e := APIRoot()
	if e != nil {
		return e
	}

	apiConf := viper.GetStringMap("api")
	if apiConf == nil {
		return errors.New("api chapter is not exist on the config file.")
	}

	var host string = apiConf["host"].(string)
	var port int = apiConf["port"].(int)

	// For CORS Header
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", host, port), c.Handler(r),
	))

	// return r.Run(fmt.Sprintf(":%d", 12345))
	// return r.Run() // listen and serve on 0.0.0.0:8080
	return nil
}
