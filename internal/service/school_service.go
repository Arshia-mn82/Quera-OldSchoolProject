package service

import (
	"OldSchool/internal/repository/models"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type SchoolRepo interface {
	Create(name string) (*models.School, error)
	List() ([]models.School, error)
	GetByID(id uint) (*models.School, error)
}

type ClassRepoForSchool interface {
	ListBySchoolID(schoolID uint) ([]models.Class, error)
}

type SchoolService struct {
	schoolRepo SchoolRepo
	classRepo  ClassRepoForSchool
}

func NewSchoolService(schoolRepo SchoolRepo, classRepo ClassRepoForSchool) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
		classRepo:  classRepo,
	}
}

func isUniqueConstraintErr(err error) bool {
	var se sqlite3.Error
	if errors.As(err, &se) {
		return se.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}

func (ss *SchoolService) Create(name string) (*models.School, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}

	created, err := ss.schoolRepo.Create(name)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return nil, ErrSchoolAlreadyExists
		}
		return nil, err
	}

	return created, nil
}

func (ss *SchoolService) List() ([]models.School, error) {
	return ss.schoolRepo.List()
}

func (ss *SchoolService) ListClasses(schoolID uint) ([]models.Class, error) {
	if schoolID == 0 {
		return nil, ErrInvalidInput
	}
	s, err := ss.schoolRepo.GetByID(schoolID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrNotFound
	}
	return ss.classRepo.ListBySchoolID(schoolID)
}
