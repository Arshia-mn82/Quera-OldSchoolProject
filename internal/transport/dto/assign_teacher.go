package dto

type AssignTeacherDTO struct {
	ClassID   uint `json:"class_id,omitempty"`
	TeacherID uint `json:"teacher_id,omitempty"`
}
