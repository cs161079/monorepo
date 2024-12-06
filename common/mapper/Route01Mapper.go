package mapper

import "github.com/cs161079/monorepo/common/models"

type Route01Mapper interface {
	GeneralRoute01(map[string]interface{}) models.Route01Oasa
	OasaToRoute01Dto(models.Route01Oasa) models.Route01
}

type route01Mapper struct {
}

func (m route01Mapper) GeneralRoute01(source map[string]interface{}) models.Route01Oasa {
	var routeDetailDto models.Route01Oasa
	internalMapper(source, &routeDetailDto)
	return routeDetailDto
}

func (m route01Mapper) OasaToRoute01Dto(source models.Route01Oasa) models.Route01 {
	var target models.Route01
	structMapper02(source, &target)
	return target
}

func NewRouteDetailMapper() Route01Mapper {
	return route01Mapper{}
}
