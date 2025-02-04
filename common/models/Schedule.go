package models

/*
***************************************************
This struct is to get data from OASA Application
*/
type Schedule struct {
	Id               int64          `json:"id" gorm:"primaryKey"`
	Sdc_Descr        string         `json:"sdc_descr"`
	Sdc_Descr_Eng    string         `json:"sdc_descr_eng"`
	Sdc_Code         int32          `json:"sdc_code"`
	Schedule_Details []Scheduletime `json:"times" gorm:"foreignKey:Sdc_Cd;references:Sdc_Code"`
}

func (Schedule) TableName() string {
	return "ScheduleMaster"
}

type ScheduleDto struct {
	Sdc_Descr     string `json:"sdc_descr"`
	Sdc_Descr_Eng string `json:"sdc_descr_eng"`
	Sdc_Code      int32  `json:"sdc_code"`
}
