package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Job struct {
	ID          string    `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Action      string    `json:"action"`
	Type        string    `json:"type"`
	Parameters  string    `json:"parameters"`

	//TBD: Convert input/output data for foreign key/association
	Input  string `json:"input"`
	Output string `json:"output"`
	AppID  string `json:"-"`
	App    App    `gorm:"ForeignID:AppID;AssociationForeignID:Id" json:"-"`
}

// BeforeCreate sets unique id to the id field
func (n *Job) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GenerateUUID())
}

// DELETE Corresponding endpoints
func (n *Job) BeforeDelete(scope *gorm.Scope) error {
	return nil
}
