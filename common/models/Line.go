package models

import "github.com/cs161079/monorepo/common/db"

type Entity interface {
}

/*
***************************************************
This struct is to get data from OASA Application
***************************************************
*/
type LineOasa struct {
	Ml_Code        int16  `json:"ml_code" oasa:"ml_code"`
	Sdc_Cd         int16  `json:"sdc_code" oasa:"sdc_code"`
	Line_Code      int32  `json:"line_code" oasa:"line_code" gorm:"index:LINE_CODE,unique"`
	Line_Id        string `json:"line_id" oasa:"line_id"`
	Line_Descr     string `json:"line_descr" oasa:"line_descr"`
	Line_Descr_Eng string `json:"line_descr_eng" oasa:"line_descr_eng"`
	Mld_master     int8   `json:"mld_master" oasa:"mld_master"`
}

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
// type Line struct {
// 	Id             int64   `json:"id" gorm:"primaryKey"`
// 	Ml_Code        int16   `json:"ml_code"`
// 	Sdc_Cd         int16   `json:"sdc_code"`
// 	Line_Code      int32   `json:"line_code" gorm:"index:Line_Code,unique"`
// 	Line_Id        string  `json:"line_id"`
// 	Line_Descr     string  `json:"line_descr"`
// 	Line_Descr_Eng string  `json:"line_descr_eng"`
// 	Mld_master     int8    `json:"mld_master"`
// 	Routes         []Route `json:"routes" gorm:"foreignKey:Ln_Code;references:line_code"`
// }

type Line struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	MLCode       int    `json:"ml_code" gorm:"column:ml_code"`
	SDCCode      int    `json:"sdc_code" gorm:"column:sdc_cd"`
	LineCode     int    `json:"line_code" gorm:"column:line_code;uniqueIndex"`
	LineID       string `json:"line_id" gorm:"column:line_id"`
	LineDescr    string `json:"line_descr" gorm:"column:line_descr"`
	LineDescrEng string `json:"line_descr_eng" gorm:"column:line_descr_eng"`
	MldMaster    int16  `json:"mld_master" gorm:"column:mld_master"`

	Routes []Route `json:"routes" gorm:"foreignKey:LnCode;references:LineCode"`
}

func (Line) TableName() string {
	return db.LINETABLE
}

type LineDto struct {
	Id             int64          `json:"id"`
	Ml_Code        int16          `json:"ml_code"`
	Sdc_Cd         int32          `json:"sdc_code"`
	Line_Code      int32          `json:"line_code"`
	Line_Id        string         `json:"line_id"`
	Line_Descr     string         `json:"line_descr"`
	Line_Descr_Eng string         `json:"line_descr_eng"`
	Mld_master     int8           `json:"mld_master"`
	Routes         []RouteDto     `json:"routes"`
	Schedule       ScheduleMaster `json:"schedule"`
}

type ComboRec struct {
	Code  int32  `json:"code"`
	Descr string `json:"descr"`
}

type LineDto01 struct {
	Ml_Code        int16  `json:"ml_code"`
	Sdc_Code       int16  `json:"sdc_code"`
	Line_Code      int32  `json:"line_code"`
	Line_Id        string `json:"line_id"`
	Line_Descr     string `json:"line_descr"`
	Line_Descr_Eng string `json:"line_descr_eng"`
}

type LineM struct {
	MLCode       int               `json:"ml_code"`
	SDCCode      int               `json:"sdc_code"`
	LineCode     int               `json:"line_code"`
	LineID       string            `json:"line_id"`
	LineDescr    string            `json:"line_descr"`
	LineDescrEng string            `json:"line_descr_eng"`
	MldMaster    int16             `json:"mld_master"`
	Routes       []RouteM          `json:"routes"`
	Schedules    []ScheduleMasterM `json:"schedules"`
}

type LineDto01M struct {
	LineCode  int         `json:"linecode"`
	LineID    string      `json:"lineid"`
	LineDescr string      `json:"linedescr"`
	Routes    []RouteDtoM `json:"routes"`
}
