package models

import "github.com/cs161079/monorepo/common/db"

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
// type Route struct {
// 	Id              int64   `json:"id" gorm:"primaryKey"`
// 	Code      		int32   `json:"route_code" gorm:"column:route_code;index:route_code_un,unique"`
// 	Ln_Code         int32   `json:"line_code"`
// 	Route_Descr     string  `json:"route_descr"`
// 	Route_Descr_eng string  `json:"route_descr_eng"`
// 	Route_Type      int8    `json:"route_type"`
// 	Route_Distance  float32 `json:"route_distance"`
// 	Route02s []Route02 `json:"stops" gorm:"foreignKey:RouteCode;references:Code"`
// }

type Route struct {
	ID            int64   `json:"id" gorm:"primaryKey"`
	RouteCode     int32   `json:"route_code" gorm:"column:route_code;uniqueIndex"`
	LnCode        int32   `json:"line_code" gorm:"column:ln_code"` // FK to Line.LineCode
	RouteDescr    string  `json:"route_descr" gorm:"column:route_descr"`
	RouteDescrEng string  `json:"route_descr_eng" gorm:"column:route_descr_eng"`
	RouteType     int8    `json:"route_type" gorm:"column:route_type"`
	RouteDistance float32 `json:"route_distance" gorm:"column:route_distance"`

	// Line     Line      `gorm:"foreignKey:LnCode;references:LineCode"`
	// Route01s []Route01 `gorm:"foreignKey:RtCode;references:RouteCode"`
	Route02s []Route02 `json:"stops" gorm:"foreignKey:RtCode;references:RouteCode"`
}

func (Route) TableName() string {
	return db.ROUTETABLE
}

/*
***************************************************
This struct is to get data from OASA Application
***************************************************
*/
type RouteOasa struct {
	Route_Code      int32   `json:"routeCode" oasa:"RouteCode"`
	Line_Code       int32   `json:"lineCode" oasa:"LineCode"`
	Route_Descr     string  `json:"routeDescr" oasa:"RouteDescr"`
	Route_Descr_Eng string  `json:"routeDescrEng" oasa:"RouteDescrEng"`
	Route_Type      int8    `json:"routeType" oasa:"RouteType"`
	Route_Distance  float32 `json:"routeDistance" oasa:"RouteDistance"`
	Stop            []Stop  `gorm:"many2many:route02;foreignKey:Route_Code;joinForeignKey:Route_Code;references:Stop_code"`
}
