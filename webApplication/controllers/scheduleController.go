package controllers

import (
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
)

type ScheduleController interface {
	getSchedule(*gin.Context)
	AddRouters(*gin.Engine)
}

type ScheduleControllerImplementation struct {
	schedSvc service.ScheduleService
}

func NewScheduleController(schedSvc service.ScheduleService) ScheduleController {
	return &ScheduleControllerImplementation{
		schedSvc: schedSvc,
	}
}

func (c ScheduleControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/schedule")
	apiGroup.GET("/details", c.getSchedule)
}

func (c ScheduleControllerImplementation) getSchedule(ctx *gin.Context) {
	start := time.Now()
	sdc_code, err := utils.StrToInt32(ctx.Query("sdc_code"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter Schedule code is not a valid number."})
		return
	}
	line_code, err := utils.StrToInt32(ctx.Query("line_code"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter Line code is not a valid number."})
		return
	}

	schedule, err := c.schedSvc.SelectByLineSdcCodeWithTimes(*line_code, *sdc_code)
	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		models.HttpResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": schedule})
}
