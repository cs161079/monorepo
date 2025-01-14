package models

import (
	"strconv"
	"strings"
)

/*
***************************************************
This struct is to get data from OASA Application
*/

// Αυτός ο τύπος είναι ένας custom τύπος για να μπορώ να μετατρέπω το JSON από τον ΟΑΣΑ
// α πεδία από String σε αριθμό.
type opswInt32 int32

// UnmarshalJSON implements custom unmarshalling logic for CustomDate
func (d *opswInt32) UnmarshalJSON(b []byte) error {
	// Parse the date from the custom format
	numStr := strings.Trim(string(b), "\"")
	if numStr == "null" {
		d = nil
	} else {
		parsedNum, err := strconv.ParseInt(numStr, 10, 32)
		if err != nil {
			return err
		}
		*d = opswInt32(int32(parsedNum))
	}
	return nil
}

type opswInt8 int8

// UnmarshalJSON implements custom unmarshalling logic for CustomDate
func (d *opswInt8) UnmarshalJSON(b []byte) error {
	// Parse the date from the custom format
	numStr := strings.Trim(string(b), "\"")
	if numStr == "null" {
		d = nil
	} else {
		parsedNum, err := strconv.ParseInt(numStr, 10, 8)
		if err != nil {
			return err
		}
		*d = opswInt8(int8(parsedNum))
	}
	return nil
}

type ScheduleDto struct {
	Sdc_Descr     string    `json:"sdc_descr"`
	Sdc_Descr_Eng string    `json:"sdc_descr_eng"`
	Sdc_Code      opswInt32 `json:"sdc_code"`
}

type Schedule struct {
	Id            int64     `json:"id" gorm:"primaryKey"`
	Sdc_Descr     string    `json:"sdc_descr"`
	Sdc_Descr_Eng string    `json:"sdc_descr_eng"`
	Sdc_Code      opswInt32 `json:"sdc_code"`
}
