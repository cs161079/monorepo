package models

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
type Route02 struct {
	Route_code int32 `json:"routeCode" gorm:"primaryKey"`
	Stop_code  int64 `json:"stopCode" gorm:"primaryKey"`
	Senu       int16 `json:"senu" gorm:"primaryKey"`
}
