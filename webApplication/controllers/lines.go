package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/webApplication/dao"
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

func (u LineControllerImplementation) GetLineInfo(ctx *gin.Context) {
	start := time.Now()
	var lineCode = ctx.Query("code")
	if lineCode == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result dao.Line
	results := u.connection.Table("line").Select("line_code, line_descr").Where("line_code=?", lineCode).Find(&result)
	if results.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"error": fmt.Sprintf("Not exists line with code=%s!", lineCode)})
		return
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	var rts []dao.Route01
	u.connection.Table("route").Select("route_code, route_descr").Where("line_code=?", lineCode).Find(&rts)
	var stps []dao.Stop01
	var rtsResult []map[string]any = make([]map[string]any, 0)
	for _, record := range rts {
		u.connection.Table("route02 routes").Select("stops.stop_code, stops.stop_descr").Joins("LEFT JOIN stop stops ON stops.stop_code=routes.stop_code").Where("routes.route_code=?", record.Route_code).Find(&stps)
		rtsResult = append(rtsResult, map[string]any{
			"code":  record.Route_code,
			"descr": record.Route_Descr,
			"stops": stps,
		})
	}
	lnInfo := map[string]any{
		"code":   result.Line_Code,
		"descr":  result.Line_Descr,
		"routes": rtsResult,
	}
	ctx.JSON(http.StatusOK, map[string]any{"data": lnInfo, "duration": time.Since(start).Seconds()})
}

func (u LineControllerImplementation) GetLineInfo02(ctx *gin.Context) {
	start := time.Now()
	var lineCode = ctx.Query("code")
	if lineCode == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	var result models.Line
	results := u.connection.Preload("Routes").Where("line_code=?", lineCode).Find(&result)
	if results.RowsAffected == 0 {
		results.Error = gorm.ErrRecordNotFound
	}
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists line with code=%s!", lineCode), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", results.Error.Error()))
		}
	}
	fmt.Printf("Query results [%d]", results.RowsAffected)
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": result})
}

func (u LineControllerImplementation) LinePreload(ctx *gin.Context) {
	var lineCode = ctx.Query("code")
	if lineCode == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}

	var line models.Line
	if err := u.connection.Preload("Routes").First(&line, "line_code=?", lineCode).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("Error occured in database selection. [%s]", err)})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"data": line})
}

func (u LineControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/line")
	apiGroup.GET("/list", u.GetLineList)
	apiGroup.GET("/info", u.GetLineInfo)
	apiGroup.GET("/info02", u.GetLineInfo02)

	apiGroup.GET("/preload", u.LinePreload)
}
