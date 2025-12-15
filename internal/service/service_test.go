package service

import (
	"OldSchool/internal/repository"
	"fmt"
	"path/filepath"
	"testing"
)

type testEnv struct {
	School *SchoolService
	Person *PersonService
	Class  *ClassService
}

func setup(t *testing.T) testEnv {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := repository.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() failed: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	// repos
	schoolRepo := repository.NewSchoolRepository(db)
	personRepo := repository.NewPersonRepositrory(db)
	classRepo := repository.NewClassRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)

	// uow
	uow := repository.NewUnitOfWork(db)

	// services
	schoolSvc := NewSchoolService(schoolRepo, classRepo)
	personSvc := NewPersonService(personRepo, classRepo, enrollRepo)
	classSvc := NewClassService(classRepo, personRepo, uow, enrollRepo)

	return testEnv{
		School: schoolSvc,
		Person: personSvc,
		Class:  classSvc,
	}
}

func TestCreateSchool_Duplicate(t *testing.T) {
	env := setup(t)

	_, err := env.School.Create("MIT")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	_, err = env.School.Create("MIT")
	fmt.Println(err)
	if err != ErrSchoolAlreadyExists {
		t.Fatalf("expected ErrSchoolAlreadyExists, got %v", err)
	}
}

func TestCreatePerson_InvalidRole(t *testing.T) {
	env := setup(t)

	_, err := env.Person.Create("Ali", "admin")
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreatePerson_OK(t *testing.T) {
	env := setup(t)

	p1, err := env.Person.Create("Teacher1", "teacher")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if p1.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}

	p2, err := env.Person.Create("Student1", "student")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if p2.ID == 0 {
		t.Fatalf("expected non-zero ID")
	}
}

func TestCreateClass_StudentCannotBeTeacher(t *testing.T) {
	env := setup(t)

	s, _ := env.School.Create("S1")
	student, _ := env.Person.Create("Stu", "student")

	_, err := env.Class.Create("Math", s.ID, student.ID)
	if err != ErrRoleMismatch {
		t.Fatalf("expected ErrRoleMismatch, got %v", err)
	}
}

func TestAddStudentToClass_Duplicate(t *testing.T) {
	env := setup(t)

	s, _ := env.School.Create("S1")
	teacher, _ := env.Person.Create("T1", "teacher")
	class, err := env.Class.Create("C1", s.ID, teacher.ID)
	if err != nil {
		t.Fatalf("create class err: %v", err)
	}

	student, _ := env.Person.Create("Stu", "student")

	// first time OK
	if err := env.Class.AddStudentToClass(student.ID, class.ID); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// second time => duplicate
	err = env.Class.AddStudentToClass(student.ID, class.ID)
	if err != ErrDuplicateEnrollment {
		t.Fatalf("expected ErrDuplicateEnrollment, got %v", err)
	}
}

func TestAddStudentToClass_DifferentSchoolRejected(t *testing.T) {
	env := setup(t)

	s1, _ := env.School.Create("S1")
	s2, _ := env.School.Create("S2")

	t1, _ := env.Person.Create("T1", "teacher")
	t2, _ := env.Person.Create("T2", "teacher")

	c1, _ := env.Class.Create("C1", s1.ID, t1.ID)
	c2, _ := env.Class.Create("C2", s2.ID, t2.ID)

	stu, _ := env.Person.Create("Stu", "student")

	// enroll in school 1
	if err := env.Class.AddStudentToClass(stu.ID, c1.ID); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// try enroll in school 2
	err := env.Class.AddStudentToClass(stu.ID, c2.ID)
	if err != ErrDifferentSchool {
		t.Fatalf("expected ErrDifferentSchool, got %v", err)
	}
}
func contains(ids []uint, target uint) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

func TestWhoAmI_TeacherAndStudent(t *testing.T) {
	env := setup(t)

	s, _ := env.School.Create("S1")
	teacher, _ := env.Person.Create("T1", "teacher")

	c1, _ := env.Class.Create("C1", s.ID, teacher.ID)
	c2, _ := env.Class.Create("C2", s.ID, teacher.ID) // ← now used

	student, _ := env.Person.Create("Stu", "student")
	_ = env.Class.AddStudentToClass(student.ID, c1.ID)

	// ---- teacher ----
	_, classIDs, err := env.Person.WhoAmI(teacher.ID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(classIDs) != 2 {
		t.Fatalf("expected 2 classes, got %d (%v)", len(classIDs), classIDs)
	}

	if !contains(classIDs, c1.ID) || !contains(classIDs, c2.ID) {
		t.Fatalf("expected classes %d and %d, got %v", c1.ID, c2.ID, classIDs)
	}

	// ---- student ----
	_, studentClassIDs, err := env.Person.WhoAmI(student.ID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(studentClassIDs) != 1 || studentClassIDs[0] != c1.ID {
		t.Fatalf("expected [%d], got %v", c1.ID, studentClassIDs)
	}
}
