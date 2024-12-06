package mapper

import models "github.com/cs161079/monorepo/common/models"

type Schedule01Mapper interface {
	DtoToSchedule01(source models.Schedule01Dto) models.Schedule01
	GeneralSchedule01(map[string]interface{}) models.Schedule01
}

type schedule01Mapper struct {
}

func (m schedule01Mapper) GeneralSchedule01(source map[string]interface{}) models.Schedule01 {
	var result models.Schedule01
	internalMapper(source, &result)
	return result
}

func (m schedule01Mapper) DtoToSchedule01(source models.Schedule01Dto) models.Schedule01 {
	var target models.Schedule01
	structMapper02(source, &target)
	return target
}

func NewSchedule01Mapper() Schedule01Mapper {
	return schedule01Mapper{}
}
