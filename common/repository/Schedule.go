package repository

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	WithTx(tx *gorm.DB) scheduleRepository
	DeleteAll() error
	SelectBySdcCodeLineCode(iLine int64, iSdc int32) (*models.Schedule, error)
	InsertScheduleMaster(input models.Schedule) error
	InsertScheduleMasterArray(input []models.Schedule) ([]models.Schedule, error)
	DeleteScheduleMaster() error

	SelectByLineSdcCodeWithTimes(int32, int32) (*models.Schedule, error)
}

type scheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(connection *gorm.DB) ScheduleRepository {
	return scheduleRepository{
		DB: connection,
	}
}

func (r scheduleRepository) DeleteAll() error {
	if err := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("1=1").Delete(&models.Schedule{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r scheduleRepository) WithTx(tx *gorm.DB) scheduleRepository {
	if tx == nil {
		logger.WARN("Database Tranction not exist.")
		return r
	}
	r.DB = tx
	return r
}

func (r scheduleRepository) SelectBySdcCodeLineCode(iLine int64, iSdc int32) (*models.Schedule, error) {
	var selectedVal models.Schedule
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("sdc_code = ? AND line_code = ?", iSdc, iLine).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return &selectedVal, nil
}

func (r scheduleRepository) InsertScheduleMaster(input models.Schedule) error {
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r scheduleRepository) DeleteScheduleMaster() error {
	if err := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("1=1").Delete(&models.Schedule{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r scheduleRepository) InsertScheduleMasterArray(input []models.Schedule) ([]models.Schedule, error) {
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Save(input)
	if res.Error != nil {
		return nil, res.Error
	}
	return input, nil
}

func (r scheduleRepository) SelectByLineSdcCodeWithTimes(lnCode int32, sdcCode int32) (*models.Schedule, error) {
	var result models.Schedule
	dbResults := r.DB.Preload("Schedule_Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("ln_code = ?", lnCode).Order("sort")
	}).Where("sdc_code = ?", sdcCode).Find(&result)

	if dbResults.Error != nil {
		return nil, fmt.Errorf("Database Error. [%s]", dbResults.Error.Error())
	}
	return &result, nil
}
