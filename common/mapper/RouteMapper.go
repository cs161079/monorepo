package mapper

import (
	"encoding/json"

	"github.com/cs161079/monorepo/common/models"
)

type RouteMapper interface {
	RouteToRouteDto(models.Route) (*models.RouteDto, error)
}

type routeMapper struct {
}

func NewRouteMapper() RouteMapper {
	return &routeMapper{}
}

func (c *routeMapper) RouteToRouteDto(orig models.Route) (*models.RouteDto, error) {
	// Transform data to JSON
	byts, err := json.Marshal(orig)
	if err != nil {
		return nil, err
	}

	// From JSON fill Dto Record
	var result models.RouteDto
	err = json.Unmarshal(byts, &result)
	if err != nil {
		return nil, err
	}

	// Transform Route02s to Array of StopDto
	result.Stops = make([]models.StopDto02, 0)
	for _, rec := range orig.Route02s {
		bytss, err := json.Marshal(rec.Stop)
		if err != nil {
			return nil, err
		}

		var stpDto models.StopDto02
		err = json.Unmarshal(bytss, &stpDto)
		if err != nil {
			return nil, err
		}
		stpDto.Senu = rec.Senu
		result.Stops = append(result.Stops, stpDto)
	}

	return &result, nil
}
