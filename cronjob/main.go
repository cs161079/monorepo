package main

import (
	"fmt"
	"os"

	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/cs161079/monorepo/cronjob/config"
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
	// ********* Κάνουμε Create το connection Με την βάση δεδομένων **************
	// dbConnection, err := db.NewOpswConnection()
	// if err != nil {
	// 	logger.ERROR(fmt.Sprintf("Database not established. [%s]", err.Error()))
	// 	return
	// }
	// ***************************************************************************
	// // ******************* Δημιουργόυμε ένα parent Context ***********************
	// var parentContext = context.Background()
	// // ******* Δημιουργούμε ένα Context με Value το connection στην Βάση *********
	// _ = context.WithValue(parentContext, db.CONNECTIONVAR, dbConnection)
	// // ***************************************************************************
	appPtr, err := config.BuildInRuntime()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	appPtr.Boot()
}
