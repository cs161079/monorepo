package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const Direction_GO = 1
const Direction_COME = 0

// CustomDate type to handle custom date formats
type opswTime time.Time

// Helper to create a new CustomTime
func NewCustomTime(hour, minute int) opswTime {
	return opswTime(time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC))
}

// UnmarshalJSON implements custom unmarshalling logic for CustomDate
func (d *opswTime) UnmarshalJSON(b []byte) error {
	// Parse the date from the custom format
	dateStr := string(b)
	if dateStr == "null" {
		d = nil
	} else {
		parsedTime, err := time.Parse(`"2006-01-02 15:04:05"`, dateStr)
		if err != nil {
			return err
		}
		*d = opswTime(parsedTime)
	}
	return nil
}

func (d opswTime) MarshalJSON() ([]byte, error) {
	// Use a custom format for JSON serialization
	ttime := time.Time(d)
	var dateStr = "null"
	if !ttime.IsZero() {
		dateStr = ttime.Format("15:04:05")
	}
	return json.Marshal(dateStr)
}

type ScheduletimeDto struct {
	Sdc_Code    opswInt32 `json:"sdc_code"`
	Line_Code   opswInt32 `json:"sde_line1"`
	Start_time1 opswTime  `json:"sde_start1"`
	End_time1   opswTime  `json:"sde_end1"`
	Start_time2 opswTime  `json:"sde_start2"`
	End_time2   opswTime  `json:"sde_end2"`
	Sort        opswInt32 `json:"sde_sort"`
}

type Scheduletime01Dto struct {
	Go   []ScheduletimeDto `json:"go"`
	Come []ScheduletimeDto `json:"come"`
}

type Scheduletime struct {
	Sdc_Cd     opswInt32 `json:"sdc_code" gorm:"primaryKey"`
	Ln_Code    opswInt32 `json:"line_code" gorm:"primaryKey"`
	Start_time opswTime  `json:"start_time" gorm:"primaryKey"`
	End_time   opswTime  `json:"end_time"`
	Sort       opswInt32 `json:"sort"`
	Direction  int8      `json:"direction" gorm:"primaryKey"`
}

func (Scheduletime) TableName() string {
	return "Scheduletime"
}

// Custom format for CustomTime
const customTimeFormat = "15:04"

// Implement the Stringer interface
func (ct opswTime) String() string {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(customTimeFormat)
}

// Implement the Scanner interface for reading from the database
func (ct *opswTime) Scan(value interface{}) error {
	if value == nil {
		*ct = opswTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = opswTime(v)
	case []uint8:
		// Parse TIME in HH:MM:SS format
		parsedTime, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return fmt.Errorf("cannot parse time: %v", err)
		}
		*ct = opswTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type %T for opswTime", value)
	}
	return nil
}

// Implement the Valuer interface for writing to the database
func (ct opswTime) Value() (driver.Value, error) {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(customTimeFormat), nil
}
