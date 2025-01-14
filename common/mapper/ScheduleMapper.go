package mapper

import (
	"encoding/json"

	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
)

type ScheduleMapper interface {
	MapDto(source any) ([]models.ScheduleDto, error)
	MapDtoToSchedule(source models.ScheduleDto) models.Schedule
}

func NewScheduleMapper() ScheduleMapper {
	return scheduleMapper{}
}

type scheduleMapper struct {
}

func (m scheduleMapper) MapDto(source any) ([]models.ScheduleDto, error) {
	var result []models.ScheduleDto
	bts, err := json.Marshal(source)
	if err != nil {
		logger.ERROR(err.Error())
		return nil, err
	}
	err = json.Unmarshal([]byte(bts), &result)
	if err != nil {
		logger.ERROR(err.Error())
		return nil, err
	}
	return result, nil
}

func (m scheduleMapper) MapDtoToSchedule(source models.ScheduleDto) models.Schedule {
	var result models.Schedule = models.Schedule{}
	structMapper02(source, &result)
	return result
}
