package controllers

import (
	"github.com/gin-gonic/gin"
)

func LcOasaRouter(masterRouter *gin.Engine) {

	oasaGroup := masterRouter.Group("/api/v1")
	{
		oasaGroup.GET("/webGetLines", WebGetLines)
		oasaGroup.GET("/webGetRoutes", WebGetRoutes)
		oasaGroup.GET("/webGetStops", WebGetStops)
		//oasaGroup.GET("/getStopNameAndXY", controller.GetStopNameAndXY)
		//oasaGroup.GET("/getLines", controller.SyncLinesApi)
		//oasaGroup.GET("/webGetMasterLines", controller.WebGetMasterlines)
	}
	//oasaSyncGroup := masterRouter.Group("/api/v1/sync")
	//{
	//	oasaSyncGroup.GET("/getLines", controller.SyncLinesApi)
	//	oasaSyncGroup.GET("/getRoutes", controller.SyncRoutesApi)
	//	oasaSyncGroup.GET("/getStops", controller.SyncStopsApi)
	//	oasaSyncGroup.GET("/getTest", controller.SyncTestApi)
	//}
}
