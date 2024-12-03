package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// convert types take an int and return a string value.
type prepareData func(*any) *any

// WebGetLines                 godoc
// @Summary      Get List of Bus Lines
// @Description  Returns all Bus Lines that is active.
// @Tags         Lines
// @Produce      json
// @Success      200  {object}  model.Line
// @Router       /webGetLines [get]
func WebGetLines(c *gin.Context) {
	// 	service.NewLineService(repository.NewLineRepository())
	// 	busLinesList, err := Busline.BuslineList01()
	// 	if err != nil {
	// 		c.IndentedJSON(http.StatusInternalServerError, struct {
	// 			Error string `json:"error"`
	// 		}{Error: err.Error()})
	// 	}
	// 	result := struct {
	// 		Lines []models.Busline `json:"lines"`
	// 	}{Lines: busLinesList}
	c.IndentedJSON(http.StatusOK, "Its ola ok")
}

// WebGetRoutes               godoc
// @Summary      Get List of Bus Routes
// @Description  Returns all Bus Routes that is active.
// @Tags         Routes
// @Param p1 query int true "Line ID"
// @Produce      json
// @Success      200  {object}  model.Route
// @Router       /webGetRoutes [get]
func WebGetRoutes(c *gin.Context) {
	// var lineCode = utils.StrToInt32(c.Query("line"))
	// routes, err := Busroute.SelectRouteByLineCode(lineCode)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, struct {
	// 		Error string `json:"error"`
	// 	}{Error: err.Error()})
	// } else {
	// 	data := struct {
	// 		Routes []models.BusRoute `json:"routes"`
	// 	}{
	// 		Routes: *routes,
	// 	}
	// 	c.IndentedJSON(http.StatusOK, data)
	// }
	c.IndentedJSON(http.StatusOK, "Its ola ok")
}

// WebGetStops               godoc
// @Summary      Get List of Bus Route's stops
// @Description  Returns all Bus stops  of route.
// @Tags         Stops
// @Param p1 query int true "Route ID"
// @Produce      json
// @Success      200  {object}  model.Stop
// @Router       /webGetStops [get]
func WebGetStops(c *gin.Context) {
	//db := database.CreateConnection()

	// var routeCode = utils.StrToInt32(c.Query("route"))
	// stopList, err := Busstop.StopList01(routeCode)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, struct {
	// 		Error string `json:"error"`
	// 	}{Error: err.Error()})
	// } else {
	// 	data := struct {
	// 		Stops []models.StopDto `json:"stops"`
	// 	}{
	// 		Stops: *stopList,
	// 	}
	// 	c.IndentedJSON(http.StatusOK, data)
	// }
	c.IndentedJSON(http.StatusOK, "its ola ok")
}
