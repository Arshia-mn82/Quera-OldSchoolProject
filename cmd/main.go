package main

import (
	"OldSchool/internal/repository"
	"OldSchool/internal/service"
	"OldSchool/internal/transport/router"
	"OldSchool/internal/transport/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := repository.InitDB("./oldSchool.db")
	if err != nil {
		log.Fatal("Cannot open the Sqlite Database")
	}

	// Repos
	schoolRepo := repository.NewSchoolRepository(db)
	personRepo := repository.NewPersonRepositrory(db)
	classRepo := repository.NewClassRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	unitOfWorkRepo := repository.NewUnitOfWork(db)

	// Services
	schoolService := service.NewSchoolService(schoolRepo, classRepo)
	personService := service.NewPersonService(personRepo, classRepo, enrollmentRepo)
	classService := service.NewClassService(classRepo, personRepo, unitOfWorkRepo, enrollmentRepo)

	// router
	router := router.NewRouter(schoolService, personService, classService)

	// server
	server := server.New(router)

	port := "8080"

	if err := server.Start(port); err != nil {
		log.Fatalf("server start failed %v", err)
	}

	log.Printf("server listening on: %s", port)

	signC := make(chan os.Signal, 1)
	signal.Notify(signC, os.Interrupt, syscall.SIGTERM)
	<-signC

	_ = server.Stop()
	log.Println("server stopped")

}
