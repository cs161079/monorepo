package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var connection *gorm.DB

type Http struct {
	Description string `json:"description,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
	StatusCode  int    `json:"statusCode"`
}

func (e Http) Error() string {
	return fmt.Sprintf("description: %s,  metadata: %s", e.Description, e.Metadata)
}

func NewHttpError(description, metadata string, statusCode int) Http {
	return Http{
		Description: description,
		Metadata:    metadata,
		StatusCode:  statusCode,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case Http:
				c.AbortWithStatusJSON(e.StatusCode, err)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
			}
		}
	}
}

func main() {
	// Δεν χρειάζεται να δώσουμε filename by Default διαβάζει .env αρχείο
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	logger.InitLogger("goSrvLogs")
	// connection, err := db.CreateConnection()
	router := gin.Default()
	router.Use(gin.Logger())
	router.Handlers = append(router.Handlers, ErrorHandler())
	logger.INFO("Router created...")

	connection, err = db.CreateConnection()
	if err != nil {
		logger.ERROR(err.Error())
	}

	// var lineService = service.NewLineService(repository.NewLineRepository(connection))

	router.GET("/getLine", func(ctx *gin.Context) {
		value, err := utils.StrToInt64(ctx.Query("id"))
		if err != nil {
			ctx.Error(err)
		}
		if *value == 0 {
			ctx.Error(NewHttpError("Δεν δόθηκε id!", "id is Empty", http.StatusBadRequest))
			return
		}
		line := models.Line{
			Id:         1234,
			Line_Descr: "Ag. Dhmhtrios - Dafni",
			Line_Id:    "131",
			Line_Code:  1234,
		}
		if *value != line.Id {
			ctx.Error(NewHttpError("Line is not exist!", fmt.Sprintf("id=%d", *value), http.StatusNotFound))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusOK, line)
	})

	router.GET("/linelist", func(ctx *gin.Context) {
		var result []models.Line = make([]models.Line, 0)
		if connection == nil {
			ctx.Error(errors.New("No Database connection was established!"))
			return
		}
		if err := connection.Table(db.LINETABLE).Where("1=1").Find(&result).Error; err != nil {
			//trans.Rollback()
			//logger.ERROR(fmt.Sprintf("An error occured on bus line list SELECT. %s", err.Error()))
			ctx.Error(err)
			return
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
