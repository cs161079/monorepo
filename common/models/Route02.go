package models

// ********* Struct for matching Stops per Route **************
// ****************** Database Entity *************************
type Route02 struct {
	Rt_code  int32 `json:"routeCode" gorm:"primaryKey;"`
	Stp_code int32 `json:"stp_code" gorm:"primaryKey"`
	Senu     int16 `json:"senu" gorm:"primaryKey"`
	Stop     Stop  `json:"stop" gorm:"foreignKey:Stop_code;references:Stp_code"`
}

func (Route02) TableName() string {
	return "Route02"
}

type Route02Dto struct {
	Stop_code  int32   `json:"stop_code"`
	Stop_descr string  `json:"stop_descr"`
	Stop_lat   float64 `json:"stop_lat"`
	Stop_lng   float64 `json:"stop_lng"`
	Senu       int16   `json:"stop_senu"`
}
