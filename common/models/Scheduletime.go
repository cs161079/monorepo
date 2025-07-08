package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cs161079/monorepo/common/db"
)

const Direction_GO = 1
const Direction_COME = 0

// ********* Struct for Schedule times **************
// *************** Database Entity ******************
// type Scheduletime struct {
// 	Sdc_Cd     int32    `json:"sdc_code" gorm:"primaryKey"`
// 	Ln_Code    int32    `json:"line_code" gorm:"primaryKey"`
// 	Start_time OpswTime `json:"start_time" gorm:"primaryKey"`
// 	End_time   OpswTime `json:"end_time"`
// 	Sort       int32    `json:"sort"`
// 	Direction  int8     `json:"direction" gorm:"primaryKey"`
// }

type ScheduleTime struct {
	LnCode    int      `json:"line_code" gorm:"column:ln_code;primaryKey"` // FK to line.line_code
	SDCCd     int      `json:"sdc_code" gorm:"column:sdc_cd;primaryKey"`   // FK to schedulemaster.sdc_code
	StartTime OpswTime `json:"start_time" gorm:"column:start_time;primaryKey"`
	EndTime   OpswTime `json:"end_time" gorm:"column:end_time"`
	Sort      int      `json:"sort" gorm:"column:sort"`
	Direction int      `json:"direction" gorm:"column:direction;primaryKey"`

	// Το έσβησα γιατ΄ιδεν ήταν υποχρεωτικό για το Preload.
	//ScheduleMaster ScheduleMaster `json:"schedule" gorm:"foreignKey:SDCCd;references:SDCCode"`
}

func (ScheduleTime) TableName() string {
	return db.SCHEDULETIMETABLE
}

// CustomDate type to handle custom date formats
type OpswTime time.Time

// Helper to create a new CustomTime
func NewCustomTime(hour, minute int) OpswTime {
	return OpswTime(time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC))
}

// UnmarshalJSON implements custom unmarshalling logic for CustomDate
func (d *OpswTime) UnmarshalJSON(b []byte) error {
	// Parse the date from the custom format
	dateStr := string(b)
	if dateStr == "null" {
		d = nil
	} else {
		parsedTime, err := time.Parse(`"15:04:05"`, dateStr)
		if err != nil {
			return err
		}
		*d = OpswTime(parsedTime)
	}
	return nil
}

func (d OpswTime) MarshalJSON() ([]byte, error) {
	// Use a custom format for JSON serialization
	ttime := time.Time(d)
	var dateStr = "null"
	if !ttime.IsZero() {
		dateStr = ttime.Format("15:04:05")
	}
	return json.Marshal(dateStr)
}

// Custom format for CustomTime
const CustomTimeFormat = "15:04"

// Implement the Stringer interface
func (ct OpswTime) String() string {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(CustomTimeFormat)
}

// Implement the Scanner interface for reading from the database
func (ct *OpswTime) Scan(value interface{}) error {
	if value == nil {
		*ct = OpswTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = OpswTime(v)
	case []uint8:
		// Parse TIME in HH:MM:SS format
		parsedTime, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return fmt.Errorf("cannot parse time: %v", err)
		}
		*ct = OpswTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type %T for opswTime", value)
	}
	return nil
}

// Implement the Valuer interface for writing to the database
func (ct OpswTime) Value() (driver.Value, error) {
	t := time.Time(ct) // Convert CustomTime to time.Time
	return t.Format(CustomTimeFormat), nil
}

type ScheduleTimeM struct {
	StartTime OpswTime `json:"start_time"`
	EndTime   OpswTime `json:"end_time"`
	Sort      int      `json:"sort"`
	Direction int      `json:"direction"`
}

// Για αυτά τα δεδομένα
// +---------+--------+--------------+
// | ln_code | sdc_cd | sdc_months   |
// +---------+--------+--------------+
// |     938 |     86 | 000001110000 |
// |     938 |     87 | 000001110000 |
// |     938 |     88 | 000001110000 |
// +---------+--------+--------------+
type ScheduleTimeDto struct {
	LnCode    int    `json:"line_code"` // FK to line.line_code
	SDCCd     int    `json:"sdc_code"`  // FK to schedulemaster.sdc_code
	SDCMonths string `json:"sdc_months"`
}
