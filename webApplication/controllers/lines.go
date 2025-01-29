package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LineController interface {
	GetLineList(*gin.Context)
	GetLineInfo(*gin.Context)
}

type LineControllerImplementation struct {
	connection *gorm.DB
	svc        service.LineService
}

func NewLineController(db *gorm.DB, svc service.LineService) LineControllerImplementation {
	return LineControllerImplementation{
		connection: db,
		svc:        svc,
	}
}

func (u LineControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/lines")
	apiGroup.GET("/list", u.GetLineList)
	apiGroup.GET("/details", u.lineDetails)
}

func (u LineControllerImplementation) GetLineList(c *gin.Context) {
	data, err := u.svc.GetLineList()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"lines": data,
	})
}

func (t LineControllerImplementation) lineDetails(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Line
	results := t.connection.Preload("Routes").Preload("Schedules.SdcDetails").Where("line_code=?", code).Find(&result)
	if results.RowsAffected == 0 {
		results.Error = gorm.ErrRecordNotFound
	}
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Line with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", results.Error.Error()))
		}
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": result})
}
