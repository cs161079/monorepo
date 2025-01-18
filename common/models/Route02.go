package models

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
type Route02 struct {
	Rt_code  int32 `json:"routeCode" gorm:"primaryKey;"`
	Stp_code int64 `json:"stopCode" gorm:"primaryKey"`
	Senu     int16 `json:"senu" gorm:"primaryKey"`
	Stop     Stop  `json:"stop" gorm:"foreignKey:Stop_code;references:Stp_code"`
}

func (Route02) TableName() string {
	return "Route02"
}
