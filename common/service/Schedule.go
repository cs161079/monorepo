package service

import "github.com/cs161079/monorepo/common/repository"

type ScheduleService interface {
}

type scheduleService struct {
	Repo repository.ScheduleRepository
}

func NewSheduleService(repo repository.ScheduleRepository) ScheduleService {
	return scheduleService{
		Repo: repo,
	}
}
