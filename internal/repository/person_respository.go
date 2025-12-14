package repository

import (
	"OldSchool/internal/repository/models"
	"errors"

	"gorm.io/gorm"
)

type PersonRepositrory struct {
	db *gorm.DB
}

func NewPersonRepositrory(db *gorm.DB) *PersonRepositrory {
	return &PersonRepositrory{db: db}
}

func (r *PersonRepositrory) Create(name string, role string) (*models.Person, error) {
	pr := &models.Person{
		Name: name,
		Role: role,
	}

	if err := r.db.Create(pr).Error; err != nil {
		return nil, err
	}

	return pr, nil
}
func (r *PersonRepositrory) GetByID(id uint) (*models.Person, error) {
	var person models.Person

	err := r.db.First(&person, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &person, nil
}

func (pr *PersonRepositrory) UpdateStudentSchoolID(studentID uint, schoolID uint) error {
	return pr.db.Model(&models.Person{}).Where("id = ?", studentID).Update("student_school_id", schoolID).Error
}
