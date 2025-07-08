package service

import (
	"sort"

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
	// routeSrv service.RouteService
}

func NewOasaService(routeRepo repository.RouteRepository, mapper mapper.OasaMapper, rest RestService) OasaService {
	return &OasaServiceImpl{
		restSrv:   rest,
		mapper:    mapper,
		routeRepo: routeRepo,
	}
}

func (c OasaServiceImpl) GetBusArrival(stop_code int32) ([]models.StopArrival, error) {

	// Request to OASA Server to ger Bus Arrivals
	response := c.restSrv.OasaRequestApi00("getStopArrivals", map[string]interface{}{"p1": stop_code})
	if response.Error != nil {
		return nil, response.Error
	}

	// var result = make([]models.StopArrival, 0)
	var structedResponse = make([]models.StopArrivalOasa, 0)
	if response.Data != nil {
		structedResponse = c.mapper.GetOasaStopArrivals(response.Data.([]interface{}))
	}

	sort.Slice(structedResponse, func(i, j int) bool {
		if structedResponse[i].Route_code == structedResponse[j].Route_code {
			return structedResponse[i].Btime2 < structedResponse[j].Btime2
		}
		return structedResponse[i].Btime2 < structedResponse[j].Btime2
	})

	var extraInfo, err = c.routeRepo.ExtraArrivalInfo(stop_code)
	if err != nil {
		return nil, err
	}

	for i, rec := range extraInfo {
		extraInfo[i].Btime2 = 999
		extraInfo[i].Last_Time = 999
		extraInfo[i].NextTime = 999

		for _, rec01 := range structedResponse {
			if rec.Route_code == rec01.Route_code {
				// if extraInfo[i].Btime2 ==
				if extraInfo[i].Btime2 == 999 {
					extraInfo[i].Veh_code = rec01.Veh_code
					extraInfo[i].Btime2 = rec01.Btime2
					var capacityInfo, err = c.routeRepo.PassengersCount(int(rec01.Veh_code), int(rec01.Route_code))
					if err != nil {
						return nil, err
					}
					extraInfo[i].Capacity = int(capacityInfo.Bus_Cap)
					extraInfo[i].Passengers = int(capacityInfo.Bus_Pass)
				} else if extraInfo[i].NextTime == 999 {
					extraInfo[i].NextTime = rec01.Btime2
				} else if extraInfo[i].Last_Time == 999 {
					extraInfo[i].Last_Time = rec01.Btime2
					break
				}
			}
		}
	}

	// // Σε αυτό το map βάζουμε βοηθητικά τις τελικές εγγραφές.
	// // Αρχικοποιήση του Map
	// var helpMap map[int32]models.StopArrival = make(map[int32]models.StopArrival)

	// // For στις γραμμές που έχω φέρει από την βάση για να εμφανίζουμε ουσιαστικά
	// // ποιες γραμμές περνάνε από αυτή τη στάση
	// for _, rec := range structedResponse {
	// 	var recRes models.StopArrival = models.StopArrival{}
	// 	// mapper για να γεμίσουμε το ένα record και τίποτα
	// 	mapper.MapStruct(rec, &recRes)
	// 	val, exists := helpMap[recRes.Route_code]
	// 	if exists {
	// 		if val.NextTime == -1 {
	// 			val.NextTime = recRes.Btime2
	// 		} else {
	// 			val.Last_Time = recRes.Btime2
	// 		}
	// 		helpMap[recRes.Route_code] = val
	// 	} else {
	// 		for _, rec01 := range extraInfo {
	// 			if rec01.Route_code == rec.Route_code {
	// 				recRes.Line_id = rec01.Line_id
	// 				recRes.Line_descr = strings.Trim(rec01.Line_descr, " ")
	// 				recRes.NextTime = -1
	// 				recRes.Last_Time = -1
	// 				helpMap[recRes.Route_code] = recRes
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// for key := range helpMap {
	// 	result = append(result, helpMap[key])
	// }

	return extraInfo, nil
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
