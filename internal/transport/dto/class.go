package dto

type CreateClassDTO struct {
	Name      string `json:"name,omitempty"`
	SchoolID  uint   `json:"school_id,omitempty"`
	TeacherID uint   `json:"teacher_id,omitempty"`
}

type AddStudentToClassDTO struct {
	StudentID uint `json:"student_id,omitempty"`
	ClassID   uint `json:"class_id,omitempty"`
}
