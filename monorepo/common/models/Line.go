package models

import "sort"

type Entity interface {
}

/*
***************************************************
This struct is to get data from OASA Application
***************************************************
*/
type LineOasa struct {
	Ml_Code        int16  `json:"masterCode" oasa:"ml_code"`
	Sdc_Code       int16  `json:"scheduleCode" oasa:"sdc_code"`
	Line_Code      int32  `json:"lineCode" oasa:"line_code" gorm:"index:LINE_CODE,unique"`
	Line_Id        string `json:"lineId" oasa:"line_id"`
	Line_Descr     string `json:"lineDescr" oasa:"line_descr"`
	Line_Descr_Eng string `json:"lineDescrEng" oasa:"line_descr_eng"`
	Mld_master     int8   `json:"mld_master" oasa:"mld_master"`
}

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
type Line struct {
	Id             int64  `json:"id" gorm:"primaryKey"`
	Ml_Code        int16  `json:"ml_code"`
	Sdc_Code       int16  `json:"sdc_code"`
	Line_Code      int32  `json:"line_code" gorm:"index:LINE_CODE,unique"`
	Line_Id        string `json:"line_id"`
	Line_Descr     string `json:"line_descr"`
	Line_Descr_Eng string `json:"line_descr_eng"`
	Mld_master     int8   `json:"mld_master"`
}

/*
*************************************************

	This struct is for different reasons

*************************************************
*/
type LineDto struct {
	Ml_Code        int64      `json:"ml_code"`
	Sdc_Code       int64      `json:"sdc_code"`
	Line_Code      int32      `json:"line_code"`
	Line_Id        string     `json:"line_id"`
	Line_Descr     string     `json:"line_descr"`
	Line_Descr_Eng string     `json:"line_descr_eng"`
	Routes         []RouteDto `json:"routes"`
	Schedules      []Schedule `json:"scheduleDay"`
}

type LineArrDto []LineDto

func (t LineArrDto) SortWithId() {
	sort.Slice(t, func(i, j int) bool {
		return t[i].Line_Id < t[j].Line_Id
	})
}

type LineDto02 struct {
	Line_id     string
	Route_Code  int32
	Route_Descr string
}
