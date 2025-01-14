package service

import (
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
	"gorm.io/gorm"
)

type ScheduleService interface {
	WithTrx(*gorm.DB) scheduleService
	DeleteAll() error
	InsertScheduleArray([]models.Schedule) ([]models.Schedule, error)
	InsertScheduleChunkArray(chunkSize int, allData []models.Schedule) error
}

type scheduleService struct {
	Repo repository.ScheduleRepository
}

func NewSheduleService(repo repository.ScheduleRepository) ScheduleService {
	return scheduleService{
		Repo: repo,
	}
}

func (s scheduleService) DeleteAll() error {
	return s.Repo.DeleteAll()
}

func (s scheduleService) WithTrx(txtHandle *gorm.DB) scheduleService {
	s.Repo = s.Repo.WithTx(txtHandle)
	return s
}

func (s scheduleService) InsertScheduleArray(allData []models.Schedule) ([]models.Schedule, error) {
	return s.Repo.InsertScheduleMasterArray(allData)
}

func (s scheduleService) InsertScheduleChunkArray(chunkSize int, allData []models.Schedule) error {
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}
	// txt := s.dbConnection.Begin()
	for {
		_, err := s.InsertScheduleArray(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.InsertScheduleArray(allData[stratIndex:])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
