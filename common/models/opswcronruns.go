package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cs161079/monorepo/common/db"
)

type OpswCronRuns struct {
	ID         int32        `json:"id" gorm:"column:id;primaryKey"`
	RUNTIME    OpswDateTime `json:"runtime" gorm:"column:runtime"`
	FINISHTIME OpswDateTime `json:"finishtime" gorm:"column:finishtime"`
	ERRORDESCR string       `json:"errorDescr" gorm:"column:errorDescr"`
}

func (OpswCronRuns) TableName() string {
	return db.OPSWCRONRUNS
}

type OpswDateTime time.Time

// Helper to create a new CustomTime
func NewDateTime() OpswDateTime {
	return OpswDateTime(time.Now())
}

// UnmarshalJSON implements custom unmarshalling logic for CustomDate
func (d *OpswDateTime) UnmarshalJSON(b []byte) error {
	// Parse the date from the custom format
	dateStr := string(b)
	if dateStr == "null" {
		d = nil
	} else {
		parsedTime, err := time.Parse(CustomDateTimeFormat, dateStr)
		if err != nil {
			return err
		}
		*d = OpswDateTime(parsedTime)
	}
	return nil
}

func (d OpswDateTime) MarshalJSON() ([]byte, error) {
	// Use a custom format for JSON serialization
	ttime := time.Time(d)
	dateStr := ttime.Format(CustomDateTimeFormat)
	fmt.Printf("%s", dateStr)
	dateStr = "null"
	if !d.IsNotSetted() {
		dateStr = ttime.Format(CustomDateTimeFormat)
	}
	return json.Marshal(dateStr)
}

// Custom format for CustomTime
const CustomDateTimeFormat = "2006-01-02 15:04:05"

// Implement the Stringer interface
func (ct OpswDateTime) String() string {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(CustomDateTimeFormat)
}

// Implement the Scanner interface for reading from the database
func (ct *OpswDateTime) Scan(value interface{}) error {
	if value == nil {
		*ct = OpswDateTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = OpswDateTime(v)
	case []uint8:
		// Parse TIME in HH:MM:SS format
		parsedTime, err := time.Parse(CustomDateTimeFormat, string(v))
		if err != nil {
			return fmt.Errorf("cannot parse time: %v", err)
		}
		*ct = OpswDateTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type %T for opswTime", value)
	}
	return nil
}

// Implement the Valuer interface for writing to the database
func (ct OpswDateTime) Value() (driver.Value, error) {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(CustomDateTimeFormat), nil
}

func (ct *OpswDateTime) IsNotSetted() bool {
	var ttime = time.Time(*ct)
	if ttime.Year() == 1 && ttime.Month() == 1 && ttime.Day() == 1 &&
		ttime.Hour() == 0 && ttime.Minute() == 0 && ttime.Second() == 0 {

		return true
	}
	return false
}
