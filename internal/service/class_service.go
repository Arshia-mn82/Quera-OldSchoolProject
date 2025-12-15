package service

import (
	"OldSchool/internal/repository"
	"OldSchool/internal/repository/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type ClassRepo interface {
	Create(name string, schoolID uint, teacherID uint) (*models.Class, error)
	GetByID(id uint) (*models.Class, error)
	ListIDsByTeacherID(teacherID uint) ([]uint, error)
	UpdateTeacher(classID, teacherID uint) error
}
type EnrollmentRepo interface {
	Exists(classID uint, studentID uint) (bool, error)
	Add(classID, studentID uint) (*models.Enrollment, error)
	ListStudentsByClassID(classID uint) ([]models.Person, error)
}

type UnitOfWork interface {
	WithinTx(fn func(r repository.Repos) error) error
}

type ClassService struct {
	classRepo      ClassRepo
	personRepo     PersonRepo
	uow            UnitOfWork
	enrollmentRepo EnrollmentRepo
}

func NewClassService(classRepo ClassRepo, personRepo PersonRepo, uow UnitOfWork, enrollmentRepo EnrollmentRepo) *ClassService {
	return &ClassService{
		uow:            uow,
		classRepo:      classRepo,
		personRepo:     personRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (cs *ClassService) Create(name string, schoolID uint, teacherID uint) (*models.Class, error) {
	name = strings.TrimSpace(name)
	if name == "" || schoolID == 0 || teacherID == 0 {
		return nil, ErrInvalidInput
	}

	teacher, err := cs.personRepo.GetByID(teacherID)
	if err != nil {
		return nil, err
	}

	if teacher == nil {
		return nil, ErrNotFound
	}

	if teacher.Role != "teacher" {
		return nil, ErrRoleMismatch
	}

	return cs.classRepo.Create(name, schoolID, teacherID)

}

func (cs *ClassService) UpdateTeacher(classID uint, teacherID uint) error {
	if classID == 0 || teacherID == 0 {
		return ErrInvalidInput
	}

	cl, err := cs.classRepo.GetByID(classID)
	if err != nil {
		return err
	}
	if cl == nil {
		return ErrNotFound
	}

	p, err := cs.personRepo.GetByID(teacherID)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrNotFound
	}
	if p.Role != "teacher" {
		return ErrRoleMismatch
	}

	if cl.TeacherID == teacherID {
		return nil
	}

	if err := cs.classRepo.UpdateTeacher(classID, teacherID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (cs *ClassService) ListStudents(classID uint) ([]models.Person, error) {
	if classID == 0 {
		return nil, ErrInvalidInput
	}

	class, err := cs.classRepo.GetByID(classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrNotFound
	}

	return cs.enrollmentRepo.ListStudentsByClassID(classID)

}

func (cs *ClassService) AddStudentToClass(studentID uint, classID uint) error {
	if studentID == 0 || classID == 0 {
		return ErrInvalidInput
	}

	student, err := cs.personRepo.GetByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return ErrNotFound
	}
	if student.Role == "teacher" {
		return ErrRoleMismatch
	}

	class, err := cs.classRepo.GetByID(classID)
	if err != nil {
		return err
	}

	if class == nil {
		return ErrNotFound
	}

	exists, err := cs.enrollmentRepo.Exists(classID, studentID)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateEnrollment
	}

	return cs.uow.WithinTx(func(r repository.Repos) error {
		st, err := r.Person.GetByID(studentID)
		if err != nil {
			return err
		}
		if st == nil {
			return ErrNotFound
		}
		if st.Role == "teacher" {
			return ErrRoleMismatch
		}

		cl, err := r.Class.GetByID(classID)
		if err != nil {
			return err
		}
		if cl == nil {
			return ErrNotFound
		}

		exists, err := r.Enrollment.Exists(classID, studentID)
		if err != nil {
			return err
		}
		if exists {
			return ErrDuplicateEnrollment
		}

		if st.StudentSchoolID == nil {
			if err := r.Person.UpdateStudentSchoolID(studentID, cl.SchoolID); err != nil {
				return err
			}
		} else if *st.StudentSchoolID != cl.SchoolID {
			return ErrDifferentSchool
		}
		_, err = r.Enrollment.Add(classID, studentID)
		return err

	})
}
