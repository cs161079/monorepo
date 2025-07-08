package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"gorm.io/gorm"
)

type busRepository struct {
	DB *gorm.DB
}

type BusRepository interface {
	SelectByBusId(int64) (*models.Bus_Capacity, error)
	SaveBusCapacity(models.Bus_Capacity) (*models.Bus_Capacity, error)
	SaveBusCapacityTest(models.Bus_Capacity) (*models.Bus_Capacity, error)
}

func NewBusRepository(connection *gorm.DB) BusRepository {
	return busRepository{
		DB: connection,
	}
}

func (r busRepository) SelectByBusId(busId int64) (*models.Bus_Capacity, error) {
	var result models.Bus_Capacity
	dbResult := r.DB.Table(db.BUSCAPACITY).Where("bus_id=?", busId).Find(&result)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &result, nil
}

func (r busRepository) SaveBusCapacity(data models.Bus_Capacity) (*models.Bus_Capacity, error) {
	//var result models.Bus_Capacity
	dbResult := r.DB.Table(db.BUSCAPACITY).Save(&data)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &data, nil
}

func (r busRepository) SaveBusCapacityTest(inputData models.Bus_Capacity) (*models.Bus_Capacity, error) {
	dbResult := r.DB.Table(db.BUSCAPACITY).Save(&inputData)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &inputData, nil
}
