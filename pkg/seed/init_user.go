package main

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"user-service/internal/config"
	"user-service/internal/domain/user/models"
	"user-service/pkg/hash"
	"user-service/pkg/logger"
)

func main() {
	db, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	InitUser(db)
	fmt.Println("Success run seeder")
}

func InitUser(db *gorm.DB) {
	password, err := hash.HashingPassword("12345678")
	if err != nil {
		logger.Error("error seeding: ", zap.Error(err))
	}
	records := []models.User{
		{
			Password:  password,
			FirstName: "Admin",
			Email:     "admin@yopmail.com",
			Role:      models.Admin,
			IsActive:  true,
		},
	}

	for _, record := range records {
		err := db.Create(&record).Error
		if err != nil {
			fmt.Printf("Error when create user: %s\n", record.FirstName)
		} else {
			fmt.Printf("Success create user: %s\n", record.FirstName)
		}
	}
}
