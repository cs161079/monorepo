package models

import (
	"github.com/cs161079/monorepo/common/db"
)

/*
***************************************************
This struct is to get data from OASA Application
*/
// type Schedule struct {
// 	Id               int64          `json:"id" gorm:"primaryKey"`
// 	Sdc_Descr        string         `json:"sdc_descr"`
// 	Sdc_Descr_Eng    string         `json:"sdc_descr_eng"`
// 	Sdc_Code         int32          `json:"sdc_code"`
// 	Sdc_days         string         `json:"sdc_days"`
// 	Sdc_months       string         `json:"sdc_months"`
// 	Schedule_Details []Scheduletime `json:"times" gorm:"foreignKey:Sdc_Cd;references:Sdc_Code"`
// }

type ScheduleMaster struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	SDCDescr    string `json:"sdc_descr" gorm:"column:sdc_descr"`
	SDCDescrEng string `json:"sdc_descr_eng" gorm:"column:sdc_descr_eng"`
	SDCCode     int32  `json:"sdc_code" gorm:"column:sdc_code;uniqueIndex"`
	SDCMonths   string `json:"sdc_months" gorm:"column:sdc_months"`
	SDCDays     string `json:"sdc_days" gorm:"column:sdc_days"`

	ScheduleTimes []ScheduleTime `json:"times" gorm:"foreignKey:SDCCd;references:SDCCode"`
}

func (ScheduleMaster) TableName() string {
	return db.SCHEDULEMASTERTABLE
}

type ScheduleDto struct {
	Sdc_Descr     string `json:"sdc_descr"`
	Sdc_Descr_Eng string `json:"sdc_descr_eng"`
	Sdc_Code      int32  `json:"sdc_code"`
}

type ScheduleMasterM struct {
	SDCDescr    string          `json:"sdc_descr"`
	SDCDescrEng string          `json:"sdc_descr_eng"`
	SDCCode     int32           `json:"sdc_code"`
	SDCMonths   string          `json:"sdc_months"`
	SDCDays     string          `json:"sdc_days"`
	Times       []ScheduleTimeM `json:"times"`
}
