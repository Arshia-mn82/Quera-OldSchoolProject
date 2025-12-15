package repository

import (
	"OldSchool/internal/repository/models"
	"errors"

	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

func (c *ClassRepository) Create(name string, schoolID uint, teacherID uint) (*models.Class, error) {
	cr := &models.Class{
		Name:      name,
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	if err := c.db.Create(cr).Error; err != nil {
		return nil, err
	}

	return cr, nil
}

func (cr *ClassRepository) GetByID(id uint) (*models.Class, error) {
	var class models.Class

	err := cr.db.First(&class, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (cr *ClassRepository) UpdateTeacher(classID, teacherID uint) error {
	tx := cr.db.Model(&models.Class{}).Where("id = ?", classID).Update("teacher_id", teacherID)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (cr *ClassRepository) ListBySchoolID(schooID uint) ([]models.Class, error) {
	var classes []models.Class
	if err := cr.db.Where("school_id = ?", schooID).Preload("Teacher").Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}
func (cr *ClassRepository) ListIDsByTeacherID(teacherID uint) ([]uint, error) {
	var ids []uint

	err := cr.db.Model(&models.Class{}).Where("teacher_id = ?", teacherID).Pluck("id", &ids).Error

	if err != nil {
		return nil, err
	}
	return ids, nil
}
