package service

import "github.com/cs161079/monorepo/common/repository"

type Schedule01Service interface {
}

type schedule01Service struct {
	Repo repository.Schedule01Repository
}
