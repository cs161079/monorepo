package mapper

type ScheduleMapper interface {
}

func NewScheduleMapper() ScheduleMapper {
	return scheduleMapper{}
}

type scheduleMapper struct {
}
