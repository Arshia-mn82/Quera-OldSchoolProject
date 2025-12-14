package main

import (
	"OldSchool/internal/repository"
	"log"
)

func main() {
	db, err := repository.InitDB("./oldSchool.db")
	if err != nil {
		log.Fatal("Cannot open the Sqlite Database")
	}
	_ = db
}
