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
}

func (t componentController) LineCombo(ctx *gin.Context) {
	start := time.Now()
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}

	var routesCb []models.ComboRec
	dbResult := t.connection.Table("Route").Select("route_code as code, route_descr as descr").Where("ln_code=?", code).Find(&routesCb)
	if dbResult.RowsAffected == 0 {
		dbResult.Error = gorm.ErrRecordNotFound
	}
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Line with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", dbResult.Error.Error()))
		}
	}

	var sdcCb []models.ComboRec
	dbResult = t.connection.Table("ScheduleMaster s").Distinct("s.sdc_code as code, s.sdc_descr as descr").Joins("LEFT JOIN ScheduleTime st ON s.sdc_code=st.sdc_cd").Where("st.ln_code=?", code).Find(&sdcCb)
	if dbResult.RowsAffected == 0 {
		dbResult.Error = gorm.ErrRecordNotFound
	}
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Line with code=%s!", code), "code": "err-001"})
			return
		} else {
			panic(fmt.Sprintln("Database Error ", dbResult.Error.Error()))
		}
	}
	var response map[string]interface{} = map[string]interface{}{"routesCb": routesCb, "sdcCb": sdcCb}
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": response})
}
