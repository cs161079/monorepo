package controllers

import (
	"net/http"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/gin-gonic/gin"
)

type RouteController interface {
	stopListByRouteCode(*gin.Context)
	AddRouters(*gin.Engine)
}

type RouteControllerImplementation struct {
	routeSvc service.RouteService
}

func NewRouteController(routeSvc service.RouteService, stopSvc service.StopService) RouteController {
	return &RouteControllerImplementation{
		routeSvc: routeSvc,
	}
}

func (u RouteControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/routes")
	apiGroup.GET("/stops", u.stopListByRouteCode)

}

func (u RouteControllerImplementation) stopListByRouteCode(ctx *gin.Context) {
	start := time.Now()
	route_code, err := utils.StrToInt32(ctx.Query("code"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": "Query Parameter code is not a valid number."})
		return
	}

	var rt models.Route
	data, err := u.routeSvc.SelectRouteWithStops(*route_code)
	if err != nil {
		// ctx.AbortWithStatusJSON(http.StatusOK, map[string]any{"error": err.Error(), "code": "err-001"})
		models.HttpResponse(ctx, err)
		return
	}
	rt = *data

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "data": rt})
}
