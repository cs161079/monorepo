package service

import (
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
	"gorm.io/gorm"
)

type ScheduleService interface {
	WithTrx(*gorm.DB) scheduleService
	DeleteAll() error
	InsertScheduleArray([]models.ScheduleMaster) ([]models.ScheduleMaster, error)
	InsertScheduleChunkArray(chunkSize int, allData []models.ScheduleMaster) error
	SelectByLineSdcCodeWithTimes(int32, int32) (*models.ScheduleMaster, error)
	SelectCurrentSchedule(int32) (*models.ScheduleMaster, error)

	ScheduleMasterList() ([]models.ScheduleMaster, error)
	// =============================================================
	//
	// Με αυτή τη διαδικασία ανασύρουμε από τη βάση δεδομένων\n
	//
	// τα Master Schedule για την γραμμή με κωδικό.\n
	//
	// @param lineCode: Κωδικός γραμμής
	//
	//
	// @return []models.ScheduleTimeDto
	//
	// @return error
	//
	// =============================================================
	ScheduleMasterDistinct(int32) ([]models.ScheduleTimeDto, error)
}

type scheduleService struct {
	repo repository.ScheduleRepository
}

func NewSheduleService(repo repository.ScheduleRepository) ScheduleService {
	return scheduleService{
		repo: repo,
	}
}

func (s scheduleService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s scheduleService) WithTrx(txtHandle *gorm.DB) scheduleService {
	s.repo = s.repo.WithTx(txtHandle)
	return s
}

func (s scheduleService) InsertScheduleArray(allData []models.ScheduleMaster) ([]models.ScheduleMaster, error) {
	return s.repo.InsertScheduleMasterArray(allData)
}

func (s scheduleService) InsertScheduleChunkArray(chunkSize int, allData []models.ScheduleMaster) error {
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

func (s scheduleService) SelectByLineSdcCodeWithTimes(lineCode int32, sdcCode int32) (*models.ScheduleMaster, error) {
	return s.repo.SelectByLineSdcCodeWithTimes(lineCode, sdcCode)
}

func (s scheduleService) SelectCurrentSchedule(linCode int32) (*models.ScheduleMaster, error) {
	return s.repo.SelectCurrentSchedule(linCode)
}

func (s scheduleService) ScheduleMasterList() ([]models.ScheduleMaster, error) {
	return s.repo.ScheduleMasterList()
}

func (s scheduleService) ScheduleMasterDistinct(lineCode int32) ([]models.ScheduleTimeDto, error) {
	return s.repo.ScheduleMasterDistinct(lineCode)
}
