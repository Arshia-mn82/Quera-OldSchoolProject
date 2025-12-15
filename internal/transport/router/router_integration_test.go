package router_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"OldSchool/internal/repository"
	"OldSchool/internal/repository/models"
	"OldSchool/internal/service"
	"OldSchool/internal/transport/protocol"
	"OldSchool/internal/transport/router"
)

func setupRouter(t *testing.T) *router.Router {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := repository.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB failed: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	schoolRepo := repository.NewSchoolRepository(db)
	personRepo := repository.NewPersonRepositrory(db)
	classRepo := repository.NewClassRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)
	uow := repository.NewUnitOfWork(db)

	schoolSvc := service.NewSchoolService(schoolRepo , classRepo)
	personSvc := service.NewPersonService(personRepo, classRepo, enrollRepo)
	classSvc := service.NewClassService(classRepo, personRepo, uow, enrollRepo)

	return router.NewRouter(schoolSvc, personSvc, classSvc)
}

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	return b
}

func TestRouter_CreateSchool_OK(t *testing.T) {
	r := setupRouter(t)

	resp := r.Handle(&protocol.Request{
		Method: router.CreateSchoolMethod,
		Data:   mustJSON(t, map[string]any{"name": "S1"}),
	})

	if !resp.Status {
		t.Fatalf("expected status=true, got %q", resp.Message)
	}

	if _, ok := resp.Data.(*models.School); !ok {
		t.Fatalf("expected *models.School, got %T", resp.Data)
	}
}

func TestRouter_CreateSchool_Duplicate(t *testing.T) {
	r := setupRouter(t)

	req := &protocol.Request{
		Method: router.CreateSchoolMethod,
		Data:   mustJSON(t, map[string]any{"name": "S1"}),
	}

	resp1 := r.Handle(req)
	if !resp1.Status {
		t.Fatalf("first create failed")
	}

	resp2 := r.Handle(req)
	if resp2.Status {
		t.Fatalf("expected duplicate to fail")
	}
}

func TestRouter_WhoAmI_TeacherAndStudent(t *testing.T) {
	r := setupRouter(t)

	// create school
	resp := r.Handle(&protocol.Request{
		Method: router.CreateSchoolMethod,
		Data:   mustJSON(t, map[string]any{"name": "S1"}),
	})
	school := resp.Data.(*models.School)

	// create teacher
	resp = r.Handle(&protocol.Request{
		Method: router.CreatePersonMethod,
		Data:   mustJSON(t, map[string]any{"name": "T1", "role": "teacher"}),
	})
	teacher := resp.Data.(*models.Person)

	// create classes C1, C2
	resp = r.Handle(&protocol.Request{
		Method: router.CreateClassMethod,
		Data: mustJSON(t, map[string]any{
			"name":       "C1",
			"school_id":  school.ID,
			"teacher_id": teacher.ID,
		}),
	})
	c1 := resp.Data.(*models.Class)

	resp = r.Handle(&protocol.Request{
		Method: router.CreateClassMethod,
		Data: mustJSON(t, map[string]any{
			"name":       "C2",
			"school_id":  school.ID,
			"teacher_id": teacher.ID,
		}),
	})

	// create student
	resp = r.Handle(&protocol.Request{
		Method: router.CreatePersonMethod,
		Data:   mustJSON(t, map[string]any{"name": "Stu", "role": "student"}),
	})
	student := resp.Data.(*models.Person)

	// enroll student in C1
	resp = r.Handle(&protocol.Request{
		Method: router.AddStudentToClassMethod,
		Data:   mustJSON(t, map[string]any{"student_id": student.ID, "class_id": c1.ID}),
	})
	if !resp.Status {
		t.Fatalf("enrollment failed: %q", resp.Message)
	}

	// whoami teacher
	resp = r.Handle(&protocol.Request{
		Method: router.WhoAmIMethod,
		Data:   mustJSON(t, map[string]any{"id": teacher.ID}),
	})
	result := resp.Data.(map[string]any)
	classIDsTeacher := result["class_ids"].([]uint)

	if len(classIDsTeacher) != 2 {
		t.Fatalf("expected 2 classes for teacher, got %d", len(classIDsTeacher))
	}

	// whoami student
	resp = r.Handle(&protocol.Request{
		Method: router.WhoAmIMethod,
		Data:   mustJSON(t, map[string]any{"id": student.ID}),
	})
	result = resp.Data.(map[string]any)
	classIDsStudent := result["class_ids"].([]uint)

	if len(classIDsStudent) != 1 || classIDsStudent[0] != c1.ID {
		t.Fatalf("unexpected student classes: %v", classIDsStudent)
	}
}
