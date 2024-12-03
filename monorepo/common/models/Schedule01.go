package models

import "time"

type rawTime []byte

func (t rawTime) Time() time.Time {
	results, err := time.Parse("15:04:05", string(t))
	if err != nil {
		panic(err.Error())
	}
	return results
}

/*
***************************************************
 */

type Schedule01Dto struct {
	Start_Time string `json:"start_time" oasa:"sde_start1" type:"time"`
	Type       int8   `json:"type"`
}

type Schedule01 struct {
	Sdc_Code   int32  `json:"sdc_code" gorm:"primaryKey"`
	Line_Code  int64  `json:"line_code" gorm:"primaryKey"`
	Start_Time string `json:"start_time" gorm:"primaryKeys"`
	Type       int8   `json:"type" gorm:"primaryKey"`
}
