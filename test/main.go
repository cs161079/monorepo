package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
)

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

	// if err := dbConnection.AutoMigrate(&User{}, &Order{}); err != nil {
	// 	logger.ERROR(fmt.Sprintf("Migration failed: %v", err))
	// }

	// Create sample data
	// dbConnection.Create(&User{
	// 	Name: "Alice",
	// 	Orders: []Order{
	// 		{Amount: 100},
	// 		{Amount: 200},
	// 	},
	// })

	// // Create sample data
	// dbConnection.Create(&User{
	// 	Name: "Nikos",
	// 	Orders: []Order{
	// 		{Amount: 15},
	// 		{Amount: 45},
	// 	},
	// })

	// Query User with Orders
	var route models.Route
	if err := dbConnection.Preload("Route02s").First(&route, "route_code=?", 1754).Error; err != nil {
		logger.ERROR(fmt.Sprintf("Query failed: %v", err))
	}

	bytes, err := json.Marshal(route)
	if err := dbConnection.Preload("Route02s").First(&route, "route_code=?", 1754).Error; err != nil {
		logger.ERROR(fmt.Sprintf("Query failed: %v", err))
	}

	log.Printf("Response: \n %s", string(bytes))

}
