package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Schedule struct {
	ID          string    `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Sched       string    `json:"schedule"`
	AppID       string    `json:"-"`
	App         App       `gorm:"ForeignID:AppID;AssociationForeignID:Id" json:"-"`
}

// BeforeCreate sets unique id to the id field
func (n *Schedule) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GenerateUUID())
}

// DELETE Corresponding endpoints
func (n *Schedule) BeforeDelete(scope *gorm.Scope) error {
	return nil
}
