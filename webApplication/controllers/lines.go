package controllers

import (
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
)

type LineController interface {
	GetLineList(*gin.Context)
	GetLineInfo(*gin.Context)
}

type LineControllerImplementation struct {
	svc          service.LineService
	routeSvc     service.RouteService
	schedService service.ScheduleService
}

func NewLineController(svc service.LineService, routeSvc service.RouteService,
	schedService service.ScheduleService) LineControllerImplementation {
	return LineControllerImplementation{
		svc:          svc,
		routeSvc:     routeSvc,
		schedService: schedService,
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
	line_code, err := utils.StrToInt32(ctx.Query("code"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter code is not a valid number."})
		return
	}

	var line *models.LineDto
	line, err = t.svc.SelectByLineCode(*line_code)
	if err != nil {
		models.HttpResponse(ctx, err)
		return
	}

	var route *models.Route
	route, err = t.routeSvc.SelectFirstRouteByLinecodeWithStops(*line_code)

	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		models.HttpResponse(ctx, err)
		return
	}

	line.Routes = append(line.Routes, *route)

	var schedule *models.Schedule

	schedule, err = t.schedService.SelectByLineSdcCodeWithTimes(line.Line_Code, line.Sdc_Cd)
	if err != nil {
		//ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		models.HttpResponse(ctx, err)
		return
	}
	line.Schedule = *schedule

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": line})
}
