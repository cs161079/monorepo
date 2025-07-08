package service

import (
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type StopService interface {
	InsertArray([]models.Stop) ([]models.Stop, error)
	InsertChunkArray(chunkSize int, allData []models.Stop) error
	DeleteAll() error
	WithTrx(*gorm.DB) stopService
	SelectByCode(int32) (*models.StopDtoBasicInfo, error)
	SelectClosestStops02(latitude float64, longtitude float64) ([]models.StopDto, error)

	SelectStops() ([]models.Stop, error)
}

type stopService struct {
	repo repository.StopRepository
}

func NewStopService(repo repository.StopRepository) StopService {
	return stopService{
		repo: repo,
	}
}

func (s stopService) WithTrx(trxHandle *gorm.DB) stopService {
	s.repo = s.repo.WithTx(trxHandle)
	return s
}

func (s stopService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s stopService) InsertArray(entityArr []models.Stop) ([]models.Stop, error) {
	return s.repo.InsertArray(entityArr)
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

func (s stopService) SelectByCode(stop_code int32) (*models.StopDtoBasicInfo, error) {
	var result models.StopDtoBasicInfo = models.StopDtoBasicInfo{}
	stopPtr, err := s.repo.SelectByCode(stop_code)
	if err != nil {
		return nil, err
	}
	mapper.MapStruct(*stopPtr, &result)
	return &result, nil
}

func (s stopService) SelectClosestStops02(latitude float64, longtitude float64) ([]models.StopDto, error) {
	return s.repo.SelectClosestStops02(latitude, longtitude)
}

func (s stopService) SelectStops() ([]models.Stop, error) {
	return s.repo.SelectAll()
}
