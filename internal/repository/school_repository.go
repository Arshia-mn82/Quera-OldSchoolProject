package repository

import (
	"OldSchool/internal/repository/models"

	"gorm.io/gorm"
)

type SchoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) Create(name string) (*models.School, error) {
	sr := &models.School{
		Name: name,
	}
	if err := r.db.Create(sr).Error; err != nil {
		return nil, err
	}
	return sr, nil	
}
