package models

import "time"

type Person struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"not null"`
	Role            string `gorm:"not null"`
	StudentSchoolID *uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
