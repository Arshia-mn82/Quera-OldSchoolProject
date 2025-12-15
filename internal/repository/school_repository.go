package repository

import (
	"OldSchool/internal/repository/models"
	"errors"

	"gorm.io/gorm"
)

type SchoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) List() ([]models.School, error) {
	var schools []models.School
	if err := r.db.Preload("Classes").Preload("Classes.Teacher").Order("id ASC").Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
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
func (r *SchoolRepository) GetByID(id uint) (*models.School, error) {
	var school models.School
	if err := r.db.First(&school, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &school, nil
}
