package main

import (
	"fmt"
	"os"

	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/cs161079/monorepo/trip-planner/config"
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
	appPtr, err := config.BuildInRuntime()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	appPtr.Boot()
}
