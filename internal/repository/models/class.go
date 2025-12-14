package models

import "time"

type Class struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	SchoolID  uint
	TeacherID *uint
	Teacher   Person `gorm:"foreignKey:TeacherID;references:ID"`
	Students []Person  `gorm:"many2many:enrollments;joinForeignKey:ClassID;joinReferences:StudentID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
