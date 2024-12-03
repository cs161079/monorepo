package service

import "github.com/cs161079/monorepo/common/repository"

type ScheduleService interface {
}

type scheduleService struct {
	Repo repository.ScheduleRepository
}
