package main

import (
	"fmt"
	"github.com/CiscoZeus/zeus-analytics/registry/api"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	"github.com/CiscoZeus/zeus-analytics/registry/template"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

// CONF_DIR is for configuration
const CONF_DIR = ".registry"

// main sets/reads configuration and execute commands
func main() {
	// TODO Is it redundant that checking current user for getting home directory ... ?
	usr, usrErr := user.Current()
	var confRoot string
	if usrErr != nil {
		// TODO Haven't tested
		confRoot = filepath.Join("~/", CONF_DIR)
	} else {
		confRoot = filepath.Join(usr.HomeDir, CONF_DIR)
	}
	viper.SetDefault("Root", confRoot)

	// If ~/.zeus-internal-gw not exist, create the directory
	if _, e := os.Stat(confRoot); os.IsNotExist(e) {
		if e2 := os.Mkdir(confRoot, os.ModeDir|0700); e2 != nil {
			fmt.Printf("Error on making %s directory: %s\n", confRoot, e2)
			os.Exit(1)
		}
	}

	// Make ~/.zeus-internal-gw/config.yml
	confFilePath := filepath.Join(confRoot, "config.yml")
	if _, e := os.Stat(confFilePath); os.IsNotExist(e) {
		confFile, e2 := os.Create(confFilePath)
		if e2 != nil {
			fmt.Printf("Could not create config file: %s\n", e2)
			os.Exit(1)
		}
		defer confFile.Close()
		confFile.Write([]byte(template.DefaultConfig()))
	}

	// Read ~/.zeus-internal-gw/config.yml
	viper.AddConfigPath(confRoot)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if e := viper.ReadInConfig(); e != nil {
		fmt.Printf("Could not read config file: %s\n", e)
		os.Exit(1)
	}

	// Migrate database model
	if e := model.Migrate(); e != nil {
		fmt.Printf("Error diring making sqlite3 database: %s\n", e)
		os.Exit(1)
	}

	if e := api.Serve(); e != nil {
		fmt.Printf("Could not launch API process: %s\n", e)
		os.Exit(1)
	}
}
