package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComponentController interface {
	LineCombo(*gin.Context)
	AddRouters(eng *gin.Engine)
}

type componentController struct {
	connection *gorm.DB
	svc        service.LineService
}

func NewComponentController(db *gorm.DB, svc service.LineService) ComponentController {
	return &componentController{
		connection: db,
		svc:        svc,
	}
}

func (u componentController) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/comp/line")
	apiGroup.GET("/cbs", u.LineCombo)
	apiGroup.GET("/alt/list", u.MasteLineCombo)
}

func (t componentController) MasteLineCombo(ctx *gin.Context) {
	start := time.Now()
	line_id := ctx.Query("line_id")
	if line_id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter Master line code is not a valid number."})
		return
	}

	comboRec, err := t.svc.AlternativeLinesList(line_id)
	if err != nil {
		models.HttpResponse(ctx, err)
		return
	}
	var response map[string]interface{} = map[string]interface{}{"altLines": comboRec}
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": response})
}

func (t componentController) LineCombo(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}
	lineCode, err := utils.StrToInt32(code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Parameter code must be a valid number [value=%s]!", code), "code": "err-001"})
	}

	var routesCb []models.ComboRec
	dbResult := t.connection.Table(db.ROUTETABLE).Select("route_code as code, route_descr as descr").Where("ln_code=?", *lineCode).Find(&routesCb)
	if dbResult.RowsAffected == 0 {
		dbResult.Error = gorm.ErrRecordNotFound
	}
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"error": fmt.Sprintf("No route were found for Line with code=%d!", *lineCode), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", dbResult.Error.Error()))
		}
	}

	var sdcCb []models.ComboRec
	dbResult = t.connection.Table(db.SCHEDULEMASTERTABLE).Distinct("schedulemaster.sdc_code as code, schedulemaster.sdc_descr as descr").
		Joins("LEFT JOIN "+db.SCHEDULETIMETABLE+" ON schedulemaster.sdc_code=scheduletime.sdc_cd").Where("scheduletime.ln_code=?", *lineCode).Find(&sdcCb)
	if dbResult.RowsAffected == 0 {
		dbResult.Error = gorm.ErrRecordNotFound
	}
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("No scheduled routes were found for Line with code=%d!", *lineCode), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", dbResult.Error.Error()))
		}
	}
	var response map[string]interface{} = map[string]interface{}{"routesCb": routesCb, "sdcCb": sdcCb}
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": response})
}
