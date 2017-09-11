package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type App struct {
	ID          string     `gorm:"primary_key" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        string     `json:"user"`
	Jobs        []Job      `json:"-"`
	Data        []Data     `json:"-"`
	Schedules   []Schedule `json:"-"`
}

// BeforeCreate sets unique id to the id field
func (n *App) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GenerateUUID())
}

// DELETE Corresponding endpoints
func (n *App) BeforeDelete(scope *gorm.Scope) error {
	// TODO Delete endpoints
	return nil
}
