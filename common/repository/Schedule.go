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

	ScheduleMasterList() ([]models.ScheduleMaster, error)
	ScheduleMasterDistinct(int32) ([]models.ScheduleTimeDto, error)
	ScheduleTimeListByLineCode(int32, int) ([]models.ScheduleTimeDto, error)
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

func (r scheduleRepository) ScheduleMasterList() ([]models.ScheduleMaster, error) {
	var dbData []models.ScheduleMaster = make([]models.ScheduleMaster, 0)
	if dbResult := r.DB.Table(db.SCHEDULEMASTERTABLE).Select("sdc_code, sdc_descr_eng, sdc_months, sdc_days").Find(&dbData); dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return dbData, nil
}

func (r scheduleRepository) ScheduleMasterDistinct(lineCode int32) ([]models.ScheduleTimeDto, error) {
	var dbData []models.ScheduleTimeDto = make([]models.ScheduleTimeDto, 0)
	dbResult := r.DB.Table(db.SCHEDULETIMETABLE).
		Distinct("scheduletime.ln_code, scheduletime.sdc_cd, schedulemaster.sdc_months, schedulemaster.sdc_days").
		Joins(fmt.Sprintf("LEFT JOIN %s ON schedulemaster.sdc_code = scheduletime.sdc_cd", db.SCHEDULEMASTERTABLE)).
		Where("scheduletime.ln_code=?", lineCode).
		Find(&dbData)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return dbData, nil
}

func (r scheduleRepository) ScheduleTimeListByLineCode(lineCode int32, direction int) ([]models.ScheduleTimeDto, error) {
	var dbData []models.ScheduleTimeDto = make([]models.ScheduleTimeDto, 0)
	dbResult := r.DB.Table(db.SCHEDULEMASTERTABLE).
		Select("scheduletime.ln_code, scheduletime.sdc_cd, schedulemaster.sdc_months, schedulemaster.sdc_days, scheduletime.start_time").
		Joins(fmt.Sprintf("LEFT JOIN %s ON schedulemaster.sdc_code = scheduletime.sdc_cd", db.SCHEDULETIMETABLE)).
		Where("scheduletime.ln_code=? AND scheduletime.direction=?", lineCode, direction).
		Order("schedulemaster.sdc_code, scheduletime.sort").
		Find(&dbData)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return dbData, nil
}
