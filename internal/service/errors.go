package service

import "errors"


var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrNotFound            = errors.New("not found")
	ErrRoleMismatch        = errors.New("role mismatch")
	ErrDuplicateEnrollment = errors.New("student already enrolled in this class")
	ErrDifferentSchool     = errors.New("student cannot enroll in multiple schools")
	ErrSchoolAlreadyExists = errors.New("school with this name already exists")
)
