package service

import "github.com/cs161079/monorepo/common/repository"

type Schedule01Service interface {
}

type schedule01Service struct {
	Repo repository.Schedule01Repository
}

func NewShedule01Service(repo repository.Schedule01Repository) Schedule01Service {
	return schedule01Service{
		Repo: repo,
	}
}
