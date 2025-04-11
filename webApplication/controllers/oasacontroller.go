package controllers

import (
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
)

type OasaNativeController interface {
	AddRouters(*gin.Engine)
}

type OasaNativeControllerImplementation struct {
	oasaSrv service.OasaService
	mapper  mapper.OasaMapper
}

func NewOasaNativeController(srv service.OasaService, mapper mapper.OasaMapper) OasaNativeController {
	return &OasaNativeControllerImplementation{
		oasaSrv: srv,
		mapper:  mapper,
	}
}

func (c OasaNativeControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/native")
	apiGroup.GET("/arrival", c.getStopArrivals)
	apiGroup.GET("/busLocation", c.getBusLocation)
}

func (c OasaNativeControllerImplementation) getStopArrivals(ctx *gin.Context) {
	start := time.Now()
	stop_codeParam := ctx.Query("code")
	stop_code, err := utils.StrToInt32(stop_codeParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter code is not a valid number."})
		return
	}

	structedResponse, err := c.oasaSrv.GetBusArrival(*stop_code)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": structedResponse})
}

func (c OasaNativeControllerImplementation) getBusLocation(ctx *gin.Context) {
	start := time.Now()
	route_codeParam := ctx.Query("code")
	route_code, err := utils.StrToInt32(route_codeParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter code is not a valid number."})
		return
	}

	structedResponse, err := c.oasaSrv.GetBusLocation(*route_code)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": structedResponse})
}
