package models

/*
***************************************************
This struct is to get data from OASA Application
*/
type ScheduleOasa struct {
	Sdc_Descr     string `json:"scheduleDescr" oasa:"sdc_descr"`
	Sdc_Descr_Eng string `json:"scheduleDescrEng" oasa:"sdc_descr_eng"`
	Sdc_Code      int32  `json:"scheduleCode" oasa:"sdc_code"`
}

type Schedule struct {
	Id            int64        `json:"id"`
	Sdc_Descr     string       `json:"sdc_descr"`
	Sdc_Descr_Eng string       `json:"sdc_descr_eng"`
	Sdc_Code      int32        `json:"sdc_code"`
	Line_Code     int64        `json:"line_code"`
	Go            []Schedule01 `oasa:"go"`
	Come          []Schedule01 `oasa:"come"`
}
