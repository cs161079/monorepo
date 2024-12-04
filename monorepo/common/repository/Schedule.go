package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	SelectBySdcCodeLineCode(iLine int64, iSdc int32) (*models.Schedule, error)
	InsertScheduleMaster(input models.Schedule)
	DeleteScheduleMaster() error
}

type scheduleRepository struct {
	DB *gorm.DB
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
