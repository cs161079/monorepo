package mapper

import models "github.com/cs161079/monorepo/common/models"

func NewOasaMapper() OasaMapper {
	return &OasaMapperImpl{}
}

type OasaMapper interface {
	GetOasaStopArrivals([]interface{}) []models.StopArrivalOasa
	GetOasaBusLocation([]interface{}) []models.BusLocation
}

type OasaMapperImpl struct {
}

func (m OasaMapperImpl) GetOasaStopArrivals(genArray []interface{}) []models.StopArrivalOasa {
	var result []models.StopArrivalOasa = make([]models.StopArrivalOasa, 0)
	for _, rec := range genArray {
		var oasaRec models.StopArrivalOasa = models.StopArrivalOasa{}
		internalMapper(rec.(map[string]interface{}), &oasaRec)
		// var stopArrival models.StopArrival = models.StopArrival{}
		// MapStruct(oasaRec, &stopArrival)
		result = append(result, oasaRec)
	}
	return result
}

func (m OasaMapperImpl) GetOasaBusLocation(genArray []interface{}) []models.BusLocation {
	var result []models.BusLocation = make([]models.BusLocation, 0)
	for _, rec := range genArray {
		var oasaRec models.BusLocation = models.BusLocation{}
		internalMapper(rec.(map[string]interface{}), &oasaRec)
		// var busLocation models.BusLocation = models.BusLocation{}
		// MapStruct(oasaRec, &busLocation)
		result = append(result, oasaRec)
	}
	return result
}
