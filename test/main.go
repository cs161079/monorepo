package main

import (
	"fmt"
	"os"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
)

func initEnviroment() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		logger.ERROR("No .env file found")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
		return
	}
	logger.InitLogger("goSyncApplication")

	dbConnection, err := db.NewOpswConnection()
	if err != nil {
		logger.ERROR(err.Error())
		return
	}

	var scheduleServ = service.NewShedule01Service(repository.NewSchedule01Repository(dbConnection))

	var inRecord models.Scheduletime = models.Scheduletime{
		Sdc_Code:   54,
		Line_Code:  1375,
		Start_time: models.NewCustomTime(13, 20),
		End_time:   models.NewCustomTime(14, 10),
		Sort:       1,
		Direction:  models.Direction_GO,
	}
	if _, err := scheduleServ.Insert(inRecord); err != nil {
		logger.ERROR(err.Error())
		return
	}
	logger.INFO("Προστέθηκε με επιτυχία.")

}
