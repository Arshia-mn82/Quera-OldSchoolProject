package dto

type SchoolClassesDTO struct {
	SchoolID uint `json:"school_id,omitempty"`
}

type ClassStudentsDTO struct {
	ClassID uint `json:"class_id,omitempty"`
}
