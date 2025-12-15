package repository

import (
	"OldSchool/internal/repository/models"
	"errors"

	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (er *EnrollmentRepository) Exists(classID uint, studentID uint) (bool, error) {
	var e models.Enrollment

	err := er.db.Where("class_id = ? AND student_id = ?", classID, studentID).First(&e).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (er *EnrollmentRepository) ListStudentsByClassID(classID uint) ([]models.Person, error) {
	var students []models.Person

	err := er.db.Table("people").Joins("JOIN enrollments ON enrollments.student_id = people.id").Where("enrollments.class_id = ?", classID).Order("people.id ASC").Scan(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (er *EnrollmentRepository) Add(classID, studentID uint) (*models.Enrollment, error) {
	e := &models.Enrollment{
		ClassID:   classID,
		StudentID: studentID,
	}

	if err := er.db.Create(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (er *EnrollmentRepository) ListClassIDsByStudentID(studentID uint) ([]uint, error) {
	var classIDs []uint

	err := er.db.Model(&models.Enrollment{}).Where("student_id = ?", studentID).Pluck("class_id", &classIDs).Error

	if err != nil {
		return nil, err
	}

	return classIDs, nil
}
