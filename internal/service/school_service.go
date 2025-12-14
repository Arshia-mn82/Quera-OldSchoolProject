package service

import (
	"OldSchool/internal/repository/models"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type SchoolRepo interface {
	Create(name string) (*models.School, error)
}

type SchoolService struct {
	schoolRepo SchoolRepo
}

func NewSchoolService(schoolRepo SchoolRepo) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
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
	}

	return created, nil
}
