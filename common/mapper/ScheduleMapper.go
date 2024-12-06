package mapper

import "github.com/cs161079/monorepo/common/models"

type ScheduleMapper interface {
	OasaToScheduleDto(source models.ScheduleOasa) models.Schedule
	DtoToSchedule(source models.Schedule) models.Schedule
	ScheduleGeneralMapper(source map[string]interface{}) models.ScheduleOasa
}

type scheduleMapper struct {
}

func (m scheduleMapper) OasaToScheduleDto(source models.ScheduleOasa) models.Schedule {
	var target models.Schedule
	structMapper02(source, &target)
	return target
}

func (m scheduleMapper) DtoToSchedule(source models.Schedule) models.Schedule {
	var target models.Schedule
	structMapper02(source, &target)
	return target
}

func (m scheduleMapper) GeneralSchedule(source map[string]interface{}) models.ScheduleOasa {
	var res models.ScheduleOasa
	internalMapper(source, &res)
	return res
}
