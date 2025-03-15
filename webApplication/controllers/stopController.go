package controllers

import (
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
)

type StopController interface {
	AddRouters(*gin.Engine)
	GetStopInfo(ctx *gin.Context)
	CloseStops(ctx *gin.Context)
}

type StopControllerImpl struct {
	stopSrv service.StopService
}

func NewStopController(srv service.StopService) StopController {
	return &StopControllerImpl{
		stopSrv: srv,
	}
}

func (u StopControllerImpl) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/stop")
	apiGroup.GET("/info", u.GetStopInfo)
	apiGroup.GET("/closeStops", u.CloseStops)
}

func (c StopControllerImpl) GetStopInfo(ctx *gin.Context) {
	start := time.Now()
	stopCodeParam := ctx.Query("code")
	stop_code, err := utils.StrToInt32(stopCodeParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter Stop code is not a valid number."})
		return
	}
	stopInfo, err := c.stopSrv.SelectByCode(*stop_code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": *stopInfo})
}

func (c StopControllerImpl) CloseStops(ctx *gin.Context) {
	start := time.Now()
	latParam := ctx.Query("lat")
	lngParam := ctx.Query("lng")
	lat := utils.StrToFloat(latParam)
	lng := utils.StrToFloat(lngParam)

	closeStops, err := c.stopSrv.SelectClosestStops02(lat, lng)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "stops": closeStops})
}
