package router

import (
	"OldSchool/internal/service"
	"OldSchool/internal/transport/dto"
	"OldSchool/internal/transport/protocol"
	"encoding/json"
)

const (
	CreateSchoolMethod      = "/school/create"
	CreateClassMethod       = "/class/create"
	CreatePersonMethod      = "/person/create"
	AddStudentToClassMethod = "/class/add/student"
	WhoAmIMethod            = "/who/am/i"
)

type Router struct {
	school *service.SchoolService
	person *service.PersonService
	class  *service.ClassService
}

func NewRouter(school *service.SchoolService, person *service.PersonService, class *service.ClassService) *Router {
	return &Router{
		school: school,
		person: person,
		class:  class,
	}
}

func ok(data any) protocol.Response {
	return protocol.Response{Status: true, Message: "ok", Data: data}
}

func badRequest(msg string) protocol.Response {
	return protocol.Response{Status: false, Message: msg, Data: nil}
}

func (r *Router) handleCreateSchoolMethod(req *protocol.Request) protocol.Response {
	var csDTO dto.CreateSchoolDTO
	if err := json.Unmarshal(req.Data, &csDTO); err != nil {
		return badRequest("invalid json for school.create")
	}
	created, err := r.school.Create(csDTO.Name)
	if err != nil {
		return fromServiceError(err)
	}
	return ok(created)
}

func (r *Router) handleCreatePersonMethod(req *protocol.Request) protocol.Response {
	var cpDTO dto.CreatePersonDTO
	if err := json.Unmarshal(req.Data, &cpDTO); err != nil {
		return badRequest("inavlid json for person.create")
	}
	created, err := r.person.Create(cpDTO.Name, cpDTO.Role)
	if err != nil {
		return fromServiceError(err)
	}
	return ok(created)
}

func (r *Router) handleCreateClassMethod(req *protocol.Request) protocol.Response {
	var ccDTo dto.CreateClassDTO
	if err := json.Unmarshal(req.Data, &ccDTo); err != nil {
		return badRequest("invalid json for class.create")
	}

	created, err := r.class.Create(ccDTo.Name, ccDTo.SchoolID, ccDTo.TeacherID)
	if err != nil {
		return fromServiceError(err)
	}
	return ok(created)
}

func (r *Router) handleAddStudentToClassMethod(req *protocol.Request) protocol.Response {
	var astcDTO dto.AddStudentToClassDTO
	if err := json.Unmarshal(req.Data, &astcDTO); err != nil {
		return badRequest("invalid json for class.add.student")
	}
	if err := r.class.AddStudentToClass(astcDTO.StudentID, astcDTO.ClassID); err != nil {
		return fromServiceError(err)
	}
	return ok(map[string]any{"status": "enrolled"})
}

func (r *Router) handleWhoAmIMethod(req *protocol.Request) protocol.Response {
	var wai dto.WhoAmIDTO
	if err := json.Unmarshal(req.Data, &wai); err != nil {
		return badRequest("inavlid json for who.am.i")
	}
	person, classIDs, err := r.person.WhoAmI(wai.ID)
	if err != nil {
		return fromServiceError(err)
	}
	return ok(map[string]any{
		"person":    person,
		"class_ids": classIDs,
	})

}
func (r *Router) Handle(req *protocol.Request) protocol.Response {
	switch req.Method {
	case CreateSchoolMethod:
		return r.handleCreateSchoolMethod(req)
	case CreatePersonMethod:
		return r.handleCreatePersonMethod(req)
	case CreateClassMethod:
		return r.handleCreateClassMethod(req)
	case AddStudentToClassMethod:
		return r.handleAddStudentToClassMethod(req)
	case WhoAmIMethod:
		return r.handleWhoAmIMethod(req)
	default:
		return protocol.Response{
			Status:  false,
			Message: "unknown method",
			Data:    nil,
		}
	}
}
