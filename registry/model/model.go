package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"path/filepath"
)

// OpenDB connects to database
func OpenDB() (*gorm.DB, error) {
	dbConf := viper.GetStringMap("db")
	if dbConf == nil {
		return nil, errors.New("db chapter is not exist on the config file.")
	}

	var driver string = dbConf["driver"].(string)
	var path string = dbConf["path"].(string)

	switch driver {
	case "sqlite3":
		root := viper.GetString("Root")
		path = filepath.Join(root, "development.sqlite3")
	case "postgres":
	default:
		fmt.Println("The driver does not found.")
	}

	return gorm.Open(driver, path)
}

// MigrateAll makes database and tables for all models
func Migrate() (err error) {
	var db *gorm.DB
	db, err = OpenDB()
	if err != nil {
		return
	}
	defer db.Close()

	db.AutoMigrate(
		&App{},
		&Job{},
		&Data{},
		&Schedule{},
	)
	return
}

func GenerateUUID() string {
	uniqueID := uuid.NewV4().String()
	return base64.StdEncoding.EncodeToString([]byte(uniqueID))
}
