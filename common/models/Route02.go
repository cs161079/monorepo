package models

import "github.com/cs161079/monorepo/common/db"

// ********* Struct for matching Stops per Route **************
// ****************** Database Entity *************************
// type Route02 struct {
// 	RouteCode   int32 `json:"routeCode" gorm:"column:rt_code;primaryKey"`
// 	Stp_code int32 `json:"stp_code" gorm:"primaryKey"`
// 	Senu     int16 `json:"senu" gorm:"primaryKey"`
// 	Stop     Stop  `json:"stop" gorm:"foreignKey:Stop_code"`
// }

type Route02 struct {
	RtCode  int32 `json:"routeCode" gorm:"column:rt_code;primaryKey"` // FK to route.route_code
	StpCode int32 `json:"stp_code" gorm:"column:stp_code;primaryKey"` // FK to stop.stop_code
	Senu    int16 `json:"senu" gorm:"column:senu;primaryKey"`

	Route Route `json:"route" gorm:"foreignKey:RtCode;references:RouteCode"`
	Stop  Stop  `json:"stop" gorm:"foreignKey:StpCode;references:StopCode"`
}

func (Route02) TableName() string {
	return db.ROUTESTOPSTABLE
}

type Route02Dto struct {
	Stop_code  int32   `json:"stop_code"`
	Stop_descr string  `json:"stop_descr"`
	Stop_lat   float64 `json:"stop_lat"`
	Stop_lng   float64 `json:"stop_lng"`
	Senu       int16   `json:"stop_senu"`
}
