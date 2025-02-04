package models

// ********* Struct for Stop **************
// ********* Database Entity **************
type Stop struct {
	Id               int64   `json:"Id" gorm:"primaryKey"`
	Stop_code        int32   `json:"Stop_code" gorm:"index:STOP_CODE_UN,unique" oasa:"StopCode"`
	Stop_id          string  `json:"stop_id" oasa:"StopID"`
	Stop_descr       string  `json:"stop_descr" oasa:"StopDescr"`
	Stop_descr_eng   string  `json:"stop_descr_eng" oasa:"StopDescrEng"`
	Stop_street      string  `json:"stop_street" oasa:"StopStreet"`
	Stop_street_eng  string  `json:"stop_street_eng" oasa:"StopStreetEng"`
	Stop_heading     int32   `json:"stop_heading" oasa:"StopHeading"`
	Stop_lat         float64 `json:"stop_lat" oasa:"StopLat"`
	Stop_lng         float64 `json:"stop_lng" oasa:"StopLng"`
	Stop_type        int8    `json:"stop_type" oasa:"StopType"`
	Stop_amea        int8    `json:"stop_amea" oasa:"StopAmea"`
	Destinations     string  `json:"destinations"`
	Destinations_Eng string  `json:"destinations_eng"`
}

func (Stop) TableName() string {
	return "stop"
}

type StopOasa struct {
	Stop_code       int32   `json:"StopCode" oasa:"StopCode"`
	Stop_id         string  `json:"stopId" oasa:"StopID"`
	Stop_descr      string  `json:"stopDescr" oasa:"StopDescr"`
	Stop_descr_eng  string  `json:"stopDescrEng" oasa:"StopDescrEng"`
	Stop_street     string  `json:"stopStreet" oasa:"StopStreet"`
	Stop_street_eng string  `json:"stopStreetEng" oasa:"StopStreetEng"`
	Stop_heading    int32   `json:"stopHeading" oasa:"StopHeading"`
	Stop_lat        float64 `json:"stopLat" oasa:"StopLat"`
	Stop_lng        float64 `json:"stopLng" oasa:"StopLng"`
	Senu            int16   `json:"stopOrder" oasa:"RouteStopOrder"`
	Stop_type       int8    `json:"stopType" oasa:"StopType"`
	Stop_amea       int8    `json:"stopAmea" oasa:"StopAmea"`
}

type StopDto struct {
	Stop_code   int32   `json:"Stop_code"`
	Stop_descr  string  `json:"stop_descr"`
	Stop_street string  `json:"stop_street"`
	Distance    float64 `json:"distance"`
}
