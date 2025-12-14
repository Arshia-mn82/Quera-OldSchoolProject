package repository

import "gorm.io/gorm"

type Repos struct {
	Person     *PersonRepositrory
	Class      *ClassRepository
	Enrollemnt *EnrollmentRepository
	School     *SchoolRepository
}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) WithinTx(fn func(r Repos) error) error {
	return uow.db.Transaction(func(tx *gorm.DB) error {
		r := Repos{
			Person:     NewPersonRepositrory(tx),
			Class:      NewClassRepository(tx),
			School:     NewSchoolRepository(tx),
			Enrollemnt: NewEnrollmentRepository(tx),
		}
		return fn(r)
	})
}
