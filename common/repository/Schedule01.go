package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type Schedule01Repository interface {
	WithTx(tx *gorm.DB) schedule01Repository
	SelectScheduleTime(lineCode int64, sdcCode int32) ([]models.ScheduleTime, error)
	SelectSchedule01ByKey(lineCode int64, sdcCode int32, tTime string, typ int) ([]models.ScheduleTime, error)
	InsertSchedule01(input models.ScheduleTime) error
	InsterSchedule01Array(input []models.ScheduleTime) ([]models.ScheduleTime, error)
	DeleteAll() error

	// -------------------- For GTFS -------------------------------------
	ScheduleTimeList(lineCode int, sdcCode int, direction int) ([]models.ScheduleTime, error)
}

type schedule01Repository struct {
	DB *gorm.DB
}

func NewSchedule01Repository(connection *gorm.DB) Schedule01Repository {
	return schedule01Repository{
		DB: connection,
	}
}

func (r schedule01Repository) WithTx(tx *gorm.DB) schedule01Repository {
	if tx == nil {
		logger.WARN("Database Tranction not exist.")
		return r
	}
	r.DB = tx
	return r
}

func (r schedule01Repository) SelectScheduleTime(lineCode int64, sdcCode int32) ([]models.ScheduleTime, error) {
	//var selectedPtr *oasaSyncModel.Busline
	var selectedVal []models.ScheduleTime
	res := r.DB.Table(db.SCHEDULETIMETABLE).Where("line_code = ? and sdc_code = ?", lineCode, sdcCode).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return selectedVal, nil
}

func (r schedule01Repository) SelectSchedule01ByKey(lineCode int64, sdcCode int32, tTime string, typ int) ([]models.ScheduleTime, error) {
	//var selectedPtr *oasaSyncModel.Busline
	var selectedVal []models.ScheduleTime
	res := r.DB.Table(db.SCHEDULETIMETABLE).Where("line_code = ? and sdc_code = ? and start_time = ? and type = ?", lineCode, sdcCode, tTime, typ).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return selectedVal, nil
}

func (r schedule01Repository) InsertSchedule01(input models.ScheduleTime) error {
	res := r.DB.Table(db.SCHEDULETIMETABLE).Create(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r schedule01Repository) InsterSchedule01Array(input []models.ScheduleTime) ([]models.ScheduleTime, error) {
	res := r.DB.Table(db.SCHEDULETIMETABLE).Save(input)
	if res.Error != nil {
		return nil, res.Error
	}
	return input, nil
}

func (r schedule01Repository) DeleteAll() error {
	if err := r.DB.Table(db.SCHEDULETIMETABLE).Where("1=1").Delete(&models.ScheduleTime{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r schedule01Repository) ScheduleTimeList(lineCode int, sdcCode int, direction int) ([]models.ScheduleTime, error) {
	var dbData []models.ScheduleTime = make([]models.ScheduleTime, 0)
	dbResults := r.DB.Table(db.SCHEDULETIMETABLE).
		Where("ln_code = ? AND sdc_cd=? AND direction = ?", lineCode, sdcCode, direction).
		Order("sort").Find(&dbData)
	if dbResults.Error != nil {
		return nil, dbResults.Error
	}
	return dbData, nil
}
