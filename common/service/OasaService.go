package service

import (
	"sort"
	"strings"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
)

type OasaService interface {
	GetBusArrival(int32) ([]models.StopArrival, error)
	GetBusLocation(int32) ([]models.BusLocation, error)
}

type OasaServiceImpl struct {
	restSrv   RestService
	mapper    mapper.OasaMapper
	routeRepo repository.RouteRepository
}

func NewOasaService(routeRepo repository.RouteRepository, mapper mapper.OasaMapper, rest RestService) OasaService {
	return &OasaServiceImpl{
		restSrv:   rest,
		mapper:    mapper,
		routeRepo: routeRepo,
	}
}

func (c OasaServiceImpl) GetBusArrival(stop_code int32) ([]models.StopArrival, error) {
	response := c.restSrv.OasaRequestApi00("getStopArrivals", map[string]interface{}{"p1": stop_code})

	if response.Error != nil {
		return nil, response.Error
	}
	var result = make([]models.StopArrival, 0)
	var structedResponse = make([]models.StopArrivalOasa, 0)
	if response.Data != nil {
		structedResponse = c.mapper.GetOasaStopArrivals(response.Data.([]interface{}))
	}
	sort.Slice(structedResponse, func(i, j int) bool {
		return structedResponse[i].Btime2 < structedResponse[j].Btime2
	})

	var extraInfo, err = c.routeRepo.ExtraArrivalInfo(stop_code)
	if err != nil {
		return nil, err
	}

	// Σε αυτό το map βάζουμε βοηθητικά τις τελικές εγγραφές.
	var helpMap map[int32]models.StopArrival = make(map[int32]models.StopArrival)
	for _, rec := range structedResponse {
		var recRes models.StopArrival = models.StopArrival{}
		mapper.MapStruct(rec, &recRes)
		val, exists := helpMap[recRes.Route_code]
		if exists {
			if val.NextTime == -1 {
				val.NextTime = recRes.Btime2
			} else {
				val.Last_Time = recRes.Btime2
			}
			helpMap[recRes.Route_code] = val
		} else {
			for _, rec01 := range extraInfo {
				if rec01.Route_code == rec.Route_code {
					recRes.Line_id = rec01.Line_id
					recRes.Line_descr = strings.Trim(rec01.Line_descr, " ")
					recRes.NextTime = -1
					recRes.Last_Time = -1
					helpMap[recRes.Route_code] = recRes
					break
				}
			}
		}
	}

	for key := range helpMap {
		result = append(result, helpMap[key])
	}

	return result, nil
}

func (s OasaServiceImpl) GetBusLocation(route_code int32) ([]models.BusLocation, error) {
	response := s.restSrv.OasaRequestApi00("getBusLocation", map[string]interface{}{"p1": route_code})

	if response.Error != nil {
		return nil, response.Error
	}
	var structedResponse = make([]models.BusLocation, 0)
	if response.Data != nil {
		structedResponse = s.mapper.GetOasaBusLocation(response.Data.([]interface{}))
	}
	return structedResponse, nil
}
