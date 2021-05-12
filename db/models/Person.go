package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	ID     int    `gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name" validate:"required"`
	Age    int    `json:"age" validate:"required"`
	Degree bool   `json:"hasDegree"`
}
