package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestController interface {
	routeInfo(*gin.Context)
	AddRoutes(*gin.Engine)
}

type TestControllerImpl struct {
	db *gorm.DB
}

func TestControllerConstructor(inDb *gorm.DB) TestController {
	return &TestControllerImpl{
		db: inDb,
	}
}

func (t TestControllerImpl) routeInfo(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Route
	results := t.db.Preload("Route01s").Where("route_code=?", code).Find(&result)
	if results.RowsAffected == 0 {
		results.Error = gorm.ErrRecordNotFound
	}
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Route with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", results.Error.Error()))
		}
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": result})
}

func (t TestControllerImpl) lineDetails(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Line
	results := t.db.Preload("Routes.Route01s").Where("line_code=?", code).Find(&result)
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

func (t TestControllerImpl) lineDetailsV1(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Line
	results := t.db.Preload("Routes").Preload("Schedules.SdcDetails").Where("line_code=?", code).Find(&result)
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

func (t TestControllerImpl) routeDetailsInfo(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Route
	results := t.db.Preload("Route01s").Preload("Route02s").Where("route_code=?", code).Find(&result)
	if results.RowsAffected == 0 {
		results.Error = gorm.ErrRecordNotFound
	}
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Route with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", results.Error.Error()))
		}
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": result})
}

func (t TestControllerImpl) routeDetailsInfoV1(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Route
	results := t.db.Preload("Route02s.Stop").Preload("Route01s").First(&result, "route_code = ?", code)
	if results.RowsAffected == 0 {
		results.Error = gorm.ErrRecordNotFound
	}
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Route with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", results.Error.Error()))
		}
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": result})
}

func (t TestControllerImpl) AddRoutes(eng *gin.Engine) {
	apiGroup := eng.Group("/test")

	apiGroup.GET("/routeInfo", t.routeInfo)
	apiGroup.GET("/lineDetails", t.lineDetails)
	apiGroup.GET("/routeDetails", t.routeDetailsInfo)
	apiGroup.GET("/v1/routeDetails", t.routeDetailsInfoV1)
	apiGroup.GET("/v1/lineDetails", t.lineDetailsV1)
}
