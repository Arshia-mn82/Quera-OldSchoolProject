package models

import "time"

type Enrollment struct {
    ClassID   uint `gorm:"primaryKey"`
    StudentID uint `gorm:"primaryKey"`
    
    CreatedAt time.Time

    Class   Class  `gorm:"foreignKey:ClassID;references:ID"`
    Student Person `gorm:"foreignKey:StudentID;references:ID"`
}
