package mapper

import models "github.com/cs161079/monorepo/common/models"

type RouteMapper interface {
	GeneralRoute(map[string]interface{}) models.RouteOasa
	OasaToRouteDto(source models.RouteOasa) models.RouteDto
	DtoToRoute(source models.RouteDto) models.Route
}

type routeMapper struct {
}

func (m routeMapper) GeneralRoute(source map[string]interface{}) models.RouteOasa {
	var busRouteOb models.RouteOasa
	internalMapper(source, &busRouteOb)

	return busRouteOb
}

func (m routeMapper) OasaToRouteDto(source models.RouteOasa) models.RouteDto {
	var target models.RouteDto
	structMapper02(source, &target)
	return target
}

func (m routeMapper) DtoToRoute(source models.RouteDto) models.Route {
	var target models.Route
	structMapper02(source, &target)
	return target
}

func NewRouteMapper() RouteMapper {
	return routeMapper{}
}
