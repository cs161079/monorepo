package service

import (
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type StopService interface {
	Post(models.Stop) (*models.Stop, error)
	InsertArray([]models.Stop) ([]models.Stop, error)
	DeleteAll() error
	WithTrx(*gorm.DB) stopService
}

type stopService struct {
	Repo repository.StopRepository
}

func NewStopService(repo repository.StopRepository) StopService {
	return stopService{
		Repo: repo,
	}
}

func (s stopService) WithTrx(trxHandle *gorm.DB) stopService {
	s.Repo = s.Repo.WithTx(trxHandle)
	return s
}

func (s stopService) Post(entity models.Stop) (*models.Stop, error) {
	routeFound, err := s.Repo.SelectByCode(entity.Stop_code)
	if err != nil {
		return nil, err
	}
	if routeFound != nil {
		return s.Repo.Update(entity)
	} else {
		return s.Repo.Insert(entity)
	}
}

func (s stopService) DeleteAll() error {
	return s.Repo.DeleteAll()
}

func (s stopService) InsertArray(entityArr []models.Stop) ([]models.Stop, error) {
	return s.Repo.InsertArray(entityArr)
}
