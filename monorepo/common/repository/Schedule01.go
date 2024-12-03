package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"gorm.io/gorm"
)

type Schedule01Repository interface {
}

type schedule01Repository struct {
	DB *gorm.DB
}

func (r schedule01Repository) SelectScheduleTime(lineCode int64, sdcCode int32) ([]models.Schedule01, error) {
	//var selectedPtr *oasaSyncModel.Busline
	var selectedVal []models.Schedule01
	res := r.DB.Table(db.SCHEDULETIMETABLE).Where("line_code = ? and sdc_code = ?", lineCode, sdcCode).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return selectedVal, nil
}

func (r schedule01Repository) SelectSchedule01ByKey(lineCode int64, sdcCode int32, tTime string, typ int) ([]models.Schedule01, error) {
	//var selectedPtr *oasaSyncModel.Busline
	var selectedVal []models.Schedule01
	res := r.DB.Table(db.SCHEDULETIMETABLE).Where("line_code = ? and sdc_code = ? and start_time = ? and type = ?", lineCode, sdcCode, tTime, typ).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return selectedVal, nil
}

func (r schedule01Repository) InsertSchedule01(input models.Schedule) error {
	res := r.DB.Table(db.SCHEDULETIMETABLE).Create(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r schedule01Repository) DeleteScheduleTime() error {
	if err := r.DB.Table(db.SCHEDULETIMETABLE).Where("1=1").Delete(&models.Schedule01{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}
