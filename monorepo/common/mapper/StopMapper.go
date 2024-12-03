package mapper

import "github.com/cs161079/monorepo/common/models"

type StopMapper interface {
	StopOasaToStop(source models.StopOasa) models.Stop
	StopMapper(source any) models.StopOasa
}

type stopMapper struct {
}

func (m stopMapper) GeneralStop(source map[string]interface{}) models.StopOasa {
	var busStopOb models.StopOasa
	internalMapper(source, &busStopOb)
	return busStopOb
}

func (m stopMapper) StopOasaToStop(source models.StopOasa) models.Stop {
	var busStop models.Stop
	structMapper02(source, &busStop)
	return busStop
}
