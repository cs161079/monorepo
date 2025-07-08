package service

import (
	"encoding/json"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
)

type BusService interface {
	SaveBusCapacity(models.BusCapacityDto) (*models.Bus_Capacity, error)
	SaveBusCapacityTest(models.BusCapacityDto) (*models.Bus_Capacity, error)
	// WithTrx(*gorm.DB) lineService
	// InsertArray([]models.Line) ([]models.Line, error)
	// InsertChunkArray(chunkSize int, allData []models.Line) error
	// DeleteAll() error
	// GetLineList() ([]models.LineDto01, error)
	// SelectByLineCode(lineCode int32) (*models.LineDto, error)

	// InsertLine(line *models.Line) (*models.Line, error)
	// PostLine(line *models.Line) (*models.Line, error)
	// PostLineArray(context.Context, []models.Line) ([]models.Line, error)
	// AlternativeLinesList(string) ([]models.ComboRec, error)

	// SearchLine(string) ([]models.Line, error)
	// GetMapper() mapper.LineMapper
}
type busService struct {
	repo repository.BusRepository
}

func NewBusService(repo repository.BusRepository) BusService {
	return busService{
		repo: repo,
	}
}

func (s busService) SaveBusCapacity(data models.BusCapacityDto) (*models.Bus_Capacity, error) {
	existingData, err := s.repo.SelectByBusId(data.Bus_Id)
	if err != nil {
		return nil, err
	}
	if existingData == nil {
		existingData = &models.Bus_Capacity{}
	}
	existingData.Bus_Id = data.Bus_Id
	existingData.Bus_Cap = data.Bus_Cap
	existingData.Bus_Pass = data.Passengers
	existingData.Date_Time = data.Date_modify
	updatedData, err := s.repo.SaveBusCapacity(*existingData)
	if err != nil {
		return nil, err
	}
	return updatedData, nil
}

func (s busService) SaveBusCapacityTest(inputData models.BusCapacityDto) (*models.Bus_Capacity, error) {
	rawData, err := json.Marshal(inputData)
	if err != nil {
		return nil, err
	}
	var preparedData models.Bus_Capacity
	json.Unmarshal(rawData, &preparedData)
	return s.repo.SaveBusCapacityTest(preparedData)
}
