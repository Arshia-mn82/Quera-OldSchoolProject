package repository

import (
	"OldSchool/internal/repository/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err == nil {
		_, _ = sqlDB.Exec("PRAGMA foreign_keys = ON")
	}

	err = db.AutoMigrate(
		&models.School{},
		&models.Person{},
		&models.Class{},
		&models.Enrollment{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil

}
