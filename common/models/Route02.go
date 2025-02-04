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
