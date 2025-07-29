package models

import "time"

type Bus_Capacity struct {
	ID        int64     `json:"id" gorm:"primaryKey;column:id"`
	Bus_Id    int64     `json:"bus_id" gorm:"column:bus_id;uniqueIndex"`
	Route_Id  int32     `json:"route_id" gorm:"column:route_id;uniqueIndex"`
	Bus_Cap   int32     `json:"bus_cap" gorm:"column:bus_cap"`
	Bus_Pass  int32     `json:"bus_pass" gorm:"column:bus_pass"`
	Date_Time time.Time `json:"date_time" gorm:"column:date_time;uniqueIndex"`
}

type BusCapacityDto struct {
	Bus_Id      int64     `json:"bus_id"`
	Route_Id    int32     `json:"route_id"`
	Bus_Cap     int32     `json:"bus_cap"`
	Passengers  int32     `json:"bus_pass"`
	Date_modify time.Time `json:"date_time"`
}

type Bus_Capacity_01 struct {
	Bus_Id     int64     `json:"bus_id" gorm:"column:bus_id;uniqueIndex"`
	Route_Id   int32     `json:"route_id" gorm:"column:route_id;uniqueIndex"`
	Sdc_Code   int32     `json:"sdc_code" gorm:"column:sdc_code;uniqueIndex"`
	Start_Time OpswTime  `json:"start_time" gorm:"column:start_time;uniqueIndex"`
	Bus_Cap    int32     `json:"bus_cap" gorm:"column:bus_cap"`
	Bus_Pass   int32     `json:"bus_pass" gorm:"column:bus_pass"`
	Date_Time  time.Time `json:"date_time" gorm:"column:date_time"`
}

type BusCapacityDto01 struct {
	Bus_Id      int64     `json:"bus_id"`
	Route_Id    int32     `json:"route_id"`
	SdcCode     int32     `json:"sdc_code"`
	StartTime   OpswTime  `json:"start_time"`
	Bus_Cap     int32     `json:"bus_cap"`
	Passengers  int32     `json:"bus_pass"`
	Date_modify time.Time `json:"date_time"`
}

type BusCapacityDt02 struct {
	Bus_Cap     int32     `json:"bus_cap"`
	Passengers  int32     `json:"bus_pass"`
	Date_modify time.Time `json:"date_time"`
}
