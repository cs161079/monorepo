package main

import (
	"fmt"
	"log"

	"github.com/cs161079/monorepo/common/db"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
)

type User struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Orders []Order `gorm:"foreignKey:UserID"` // Correctly establishes the relationship
}

type Order struct {
	ID     uint `gorm:"primaryKey"`
	Amount float64
	UserID uint // Foreign key that references User
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
		return
	}
	logger.InitLogger("goTestApplication")

	dbConnection, err := db.NewOpswConnection()
	if err != nil {
		logger.ERROR(err.Error())
		return
	}

	if err := dbConnection.AutoMigrate(&User{}, &Order{}); err != nil {
		logger.ERROR(fmt.Sprintf("Migration failed: %v", err))
	}

	// Create sample data
	dbConnection.Create(&User{
		Name: "Alice",
		Orders: []Order{
			{Amount: 100},
			{Amount: 200},
		},
	})

	// Create sample data
	dbConnection.Create(&User{
		Name: "Nikos",
		Orders: []Order{
			{Amount: 15},
			{Amount: 45},
		},
	})

	// Query User with Orders
	var user User
	if err := dbConnection.Preload("Orders").First(&user, "id=?", 1).Error; err != nil {
		logger.ERROR(fmt.Sprintf("Query failed: %v", err))
	}

	log.Printf("User: %+v", user)

}
