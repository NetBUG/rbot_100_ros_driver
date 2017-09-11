package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Data struct {
	ID          string    `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Type        string    `json:"type"`
	Store       string    `json:"store"`
	Parameters  string    `json:"parameters"`
	AppID       string    `json:"-"`
	App         App       `gorm:"ForeignID:AppID;AssociationForeignID:Id" json:"-"`
}

// BeforeCreate sets unique id to the id field
func (n *Data) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GenerateUUID())
}

// DELETE Corresponding endpoints
func (n *Data) BeforeDelete(scope *gorm.Scope) error {
	return nil
}
