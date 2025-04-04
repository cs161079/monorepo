package service

import (
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
	"gorm.io/gorm"
)

type Schedule01Service interface {
	WithTrx(*gorm.DB) schedule01Service
	DeleteAll() error
	InsertArray(allData []models.ScheduleTime) ([]models.ScheduleTime, error)
	InsertSchedule01ChunkArray(chunkSize int, allData []models.ScheduleTime) error
}

type schedule01Service struct {
	repo repository.Schedule01Repository
}

func NewShedule01Service(repo repository.Schedule01Repository) Schedule01Service {
	return schedule01Service{
		repo: repo,
	}
}

func (s schedule01Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s schedule01Service) WithTrx(txtHandle *gorm.DB) schedule01Service {
	s.repo = s.repo.WithTx(txtHandle)
	return s
}

func (s schedule01Service) InsertArray(allData []models.ScheduleTime) ([]models.ScheduleTime, error) {
	return s.repo.InsterSchedule01Array(allData)
}

func (s schedule01Service) InsertSchedule01ChunkArray(chunkSize int, allData []models.ScheduleTime) error {
	// var maxSize = 1000
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}
	// txt := s.dbConnection.Begin()
	for {
		_, err := s.InsertArray(allData[stratIndex:endIndex])
		if err != nil {
			// txt.Rollback()
			//logger.ERROR(fmt.Sprintf("Σφάλμα κατά την προσθήκη των γραμμών από %d έως %d.", stratIndex, endIndex-1))
			return err
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι γραμμές από %d έως %d.", stratIndex, endIndex-1))
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			//logger.INFO("Η εισαγωγή γραμμών ολοκληρώθηκε.")
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.InsertArray(allData[stratIndex:])
			if err != nil {
				//txt.Rollback()
				//logger.ERROR(fmt.Sprintf("Σφάλμα κατά την προσθήκη των γραμμών από %d έως Τέλος.", stratIndex))
				return err
			}
			break
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι γραμμές από %d έως %d.", stratIndex, endIndex-1))
	}
	return nil
}
