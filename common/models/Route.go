package models

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
	Stop            []Stop  `gorm:"many2many:route02;foreignKey:Route_Code;joinForeignKey:Route_Code;References:Stop_code"`
}

/*
******************************************
Struct for Bus Lines Entities for database
******************************************
*/
type Route struct {
	Id              int64   `json:"Id" gorm:"primaryKey"`
	Route_Code      int32   `json:"route_code" gorm:"index:route_code_un,unique"`
	Ln_Code         int32   `json:"line_id"`
	Route_Descr     string  `json:"route_descr"`
	Route_Descr_eng string  `json:"route_descr_eng"`
	Route_Type      int8    `json:"route_type"`
	Route_Distance  float32 `json:"route_distance"`
	// Route01s        []Route01 `json:"route_details" gorm:"foreignKey:Rt_code;references:Route_Code"`
	// Route02s        []Route02 `json:"stops" gorm:"foreignKey:Rt_code;references:Route_Code"`
}

func (Route) TableName() string {
	return "Route"
}

/*
	*************************************************
	       This struct is for different reasons
	*************************************************
*/

type RouteDto struct {
	Id              int64     `json:"Id"`
	Route_Code      int32     `json:"route_code"`
	Line_Code       int32     `json:"line_code"`
	Route_Descr     string    `json:"route_descr"`
	Route_Descr_Eng string    `json:"route_descr_eng"`
	Route_Type      int8      `json:"route_type"`
	Route_Distance  float32   `json:"route_distance"`
	Stops           []Stop    `json:"stops"`
	RouteDetails    []Route01 `json:"routeDetails"`
}

//**********************************************************
