package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
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
	var code = ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Line code must have a value."})
		return
	}

	var line *models.LineDto
	line, err := t.svc.SelectByLineCode(code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": fmt.Sprintf("Not exists Line with code=%s!", code), "code": "err-001"})
		return
	}

	var route *models.Route
	route, err = t.routeSvc.SelectFirstRouteByLinecodeWithStops(code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		return
	}

	line.Routes = append(line.Routes, *route)

	var schedule *models.Schedule

	schedule, err = t.schedService.SelectByLineSdcCodeWithTimes(code, line.Sdc_Cd)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		return
	}
	line.Schedule = *schedule

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": line})
}
