package repository

import (
	"fmt"
	"time"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	WithTx(tx *gorm.DB) scheduleRepository
	DeleteAll() error
	SelectBySdcCodeLineCode(iLine int64, iSdc int32) (*models.ScheduleMaster, error)
	InsertScheduleMaster(input models.ScheduleMaster) error
	InsertScheduleMasterArray(input []models.ScheduleMaster) ([]models.ScheduleMaster, error)
	DeleteScheduleMaster() error

	SelectByLineSdcCodeWithTimes(int32, int32) (*models.ScheduleMaster, error)
	SelectCurrentSchedule(int32) (*models.ScheduleMaster, error)
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
	if err := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("1=1").Delete(&models.ScheduleMaster{}).Error; err != nil {
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

func (r scheduleRepository) SelectBySdcCodeLineCode(iLine int64, iSdc int32) (*models.ScheduleMaster, error) {
	var selectedVal models.ScheduleMaster
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("sdc_code = ? AND line_code = ?", iSdc, iLine).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return &selectedVal, nil
}

func (r scheduleRepository) InsertScheduleMaster(input models.ScheduleMaster) error {
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r scheduleRepository) DeleteScheduleMaster() error {
	if err := r.DB.Table(db.SCHEDULEMASTERTABLE).Where("1=1").Delete(&models.ScheduleMaster{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r scheduleRepository) InsertScheduleMasterArray(input []models.ScheduleMaster) ([]models.ScheduleMaster, error) {
	res := r.DB.Table(db.SCHEDULEMASTERTABLE).Save(input)
	if res.Error != nil {
		return nil, res.Error
	}
	return input, nil
}

func (r scheduleRepository) SelectByLineSdcCodeWithTimes(lineCode int32, sdcCode int32) (*models.ScheduleMaster, error) {
	var result models.ScheduleMaster
	dbResults := r.DB.Preload("ScheduleTimes", func(db *gorm.DB) *gorm.DB {
		return db.Where("ln_code = ?", lineCode).Order("sort")
	}).Where("sdc_code=?", sdcCode).First(&result)

	if dbResults.Error != nil {
		return nil, fmt.Errorf("Database Error. [%s]", dbResults.Error.Error())
	}

	return &result, nil
}

func (r scheduleRepository) SelectCurrentSchedule(lineCode int32) (*models.ScheduleMaster, error) {
	var result models.ScheduleMaster
	var hlpArr []models.ScheduleMaster
	dbResults := r.DB.Preload("Schedule_Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("ln_code = ?", lineCode).Order("sort")
	}).Find(&hlpArr)

	if dbResults.Error != nil {
		return nil, fmt.Errorf("Database Error. [%s]", dbResults.Error.Error())
	}

	currentMonth := int(time.Now().Month())
	currentDay := time.Now().Weekday()

	for _, rec := range hlpArr {
		if len(rec.ScheduleTimes) != 0 && string(rec.SDCDays[currentDay]) == "1" && string(rec.SDCMonths[currentMonth-1]) == "1" {
			result = rec
			return &result, nil
		}
	}

	return &result, nil
}
