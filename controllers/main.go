package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var connection *gorm.DB

func main() {
	// Δεν χρειάζεται να δώσουμε filename by Default διαβάζει .env αρχείο
	err := godotenv.Load()

	if err != nil {
		fmt.Printf(err.Error())
	}

	logger.InitLogger("goSrvLogs")
	// connection, err := db.CreateConnection()
	router := gin.Default()
	logger.INFO("Router created...")

	connection, err = db.CreateConnection()
	if err != nil {
		logger.ERROR(err.Error())
	}

	// var lineService = service.NewLineService(repository.NewLineRepository(connection))

	router.GET("/linelist", func(ctx *gin.Context) {
		var result []models.Line = make([]models.Line, 0)
		if err := connection.Table(db.LINETABLE).Where("1=1").Find(&result).Error; err != nil {
			//trans.Rollback()
			logger.ERROR(fmt.Sprintf("An error occured on bus line list SELECT. %s", err.Error()))
		}
		var finalResult = map[string]interface{}{
			"lines": result,
		}
		ctx.IndentedJSON(http.StatusOK, finalResult)
	})

	router.GET("/test", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "This is a test API, ok")

	})

	router.GET("/test02", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "This is a test1 API, ok")
	})

	logger.INFO("GIN Server started...")
	var serverPort = os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8080"
	}
	router.Run(":" + serverPort)

}
