package models

import "time"

type School struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null;uniqueIndex"`
	Classes   []Class `gorm:"foreignKey:SchoolID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
