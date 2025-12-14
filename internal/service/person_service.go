package service

import (
	"OldSchool/internal/repository/models"
	"strings"
)

type PersonRepo interface {
	Create(name string, role string) (*models.Person, error)
	GetByID(id uint) (*models.Person, error)
}

type ClassRepoWhoAmI interface {
	ListIDsByTeacherID(teacheriD uint) ([]uint, error)
}

type EnrollmentRepoForWhoAmI interface {
	ListClassIDsByStudentID(studentID uint) ([]uint, error)
}

type PersonService struct {
	personRepo     PersonRepo
	classRepo      ClassRepoWhoAmI
	enrollmentRepo EnrollmentRepoForWhoAmI
}

func NewPersonService(personRepo PersonRepo, classRepo ClassRepoWhoAmI, enrollmentRepo EnrollmentRepoForWhoAmI) *PersonService {
	return &PersonService{
		personRepo:     personRepo,
		classRepo:      classRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (pr *PersonService) Create(name, role string) (*models.Person, error) {
	name = strings.TrimSpace(name)
	role = strings.TrimSpace(role)

	if name == "" {
		return nil, ErrInvalidInput
	}

	if role != "teacher" && role != "student" {
		return nil, ErrInvalidInput
	}

	created, err := pr.personRepo.Create(name, role)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (pr *PersonService) WhoAmI(personID uint) (*models.Person, []uint, error) {
	p, err := pr.personRepo.GetByID(personID)
	if err != nil {
		return nil, nil, err
	}

	if p == nil {
		return nil, nil, ErrNotFound
	}

	switch p.Role {
	case "student":
		classIDs, err := pr.enrollmentRepo.ListClassIDsByStudentID(p.ID)
		if err != nil {
			return nil, nil, err
		}
		return p, classIDs, nil
	case "teacher":
		classIDs, err := pr.classRepo.ListIDsByTeacherID(p.ID)
		if err != nil {
			return nil, nil, err
		}
		return p, classIDs, nil
	default:
		return nil, nil, ErrRoleMismatch
	}
}
