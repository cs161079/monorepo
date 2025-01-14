package service

import (
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type StopService interface {
	Post(models.Stop) (*models.Stop, error)
	InsertArray([]models.Stop) ([]models.Stop, error)
	InsertChunkArray(chunkSize int, allData []models.Stop) error
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

func (s stopService) InsertChunkArray(chunkSize int, allData []models.Stop) error {
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}

	for {
		_, err := s.InsertArray(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι λεπτομερειες διαδρομών από %d έως %d.", stratIndex, endIndex-1))
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			//logger.INFO("Η εισαγωγή λεπτομερειών διαδρομών ολοκληρώθηκε.")
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.InsertArray(allData[stratIndex:])
			if err != nil {
				//logger.ERROR(fmt.Sprintf("Σφάλμα κατά την προσθήκη λεπτομερειών διαδρομών από %d έως τέλος.", stratIndex))
				//txt.Rollback()
				return err
			}
			break
		}
	}
	return nil
}
