package models

import "github.com/cs161079/monorepo/common/db"

// ********* Struct for Stop **************
// ********* Database Entity **************
// type Stop struct {
// 	Id               int64   `json:"Id" gorm:"primaryKey"`
// 	Stop_code        int32   `json:"stop_code" gorm:"index:STOP_CODE_UN,unique"`
// 	Stop_id          string  `json:"stop_id"`
// 	Stop_descr       string  `json:"stop_descr"`
// 	Stop_descr_eng   string  `json:"stop_descr_eng"`
// 	Stop_street      string  `json:"stop_street"`
// 	Stop_street_eng  string  `json:"stop_street_eng"`
// 	Stop_heading     int32   `json:"stop_heading"`
// 	Stop_lat         float64 `json:"stop_lat"`
// 	Stop_lng         float64 `json:"stop_lng"`
// 	Stop_type        int8    `json:"stop_type"`
// 	Stop_amea        int8    `json:"stop_amea"`
// 	Destinations     string  `json:"destinations"`
// 	Destinations_Eng string  `json:"destinations_eng"`
// }

type Stop struct {
	ID              int     `json:"id" gorm:"primaryKey"`
	StopCode        int32   `json:"stop_code" gorm:"column:stop_code;uniqueIndex"`
	StopID          string  `json:"stop_id" gorm:"column:stop_id"`
	StopDescr       string  `json:"stop_descr" gorm:"column:stop_descr"`
	StopDescrEng    string  `json:"stop_descr_eng" gorm:"column:stop_descr_eng"`
	StopStreet      string  `json:"stop_street" gorm:"column:stop_street"`
	StopStreetEng   string  `json:"stop_street_eng" gorm:"column:stop_street_eng"`
	StopHeading     int32   `json:"stop_heading" gorm:"column:stop_heading"`
	StopLat         float64 `json:"stop_lat" gorm:"column:stop_lat"`
	StopLng         float64 `json:"stop_lng" gorm:"column:stop_lng"`
	StopType        int8    `json:"stop_type" gorm:"column:stop_type"`
	StopAmea        int8    `json:"stop_amea" gorm:"column:stop_amea"`
	Destinations    string  `json:"destinations" gorm:"column:destinations"`
	DestinationsEng string  `json:"destinations_eng" gorm:"column:destinations_eng"`

	RouteStops []Route02 `json:"route_stops" gorm:"foreignKey:StpCode;references:StopCode"`
}

func (Stop) TableName() string {
	return db.STOPTABLE
}

type StopDto02 struct {
	StopCode     int32  `json:"stop_code"`
	StopID       string `json:"stop_id"`
	StopDescr    string `json:"stop_descr"`
	StopDescrEng string `json:"stop_descr_eng"`
	Senu         int16  `json:"senu"`
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
	Stop_code   int32   `json:"stop_code"`
	Stop_descr  string  `json:"stop_descr"`
	Stop_street string  `json:"stop_street"`
	Stop_lat    float64 `json:"stop_lat"`
	Stop_lng    float64 `json:"stop_lng"`
	Distance    float64 `json:"distance"`
}

// Αυτό μπορεί να το αλλάξω.
// Το ήθελα για την σελίδα της στάσης
type StopDtoBasicInfo struct {
	StopID     string `json:"stop_id"`
	Stop_code  int32  `json:"stop_code"`
	Stop_descr string `json:"stop_descr"`
}

type StopM struct {
	StopCode        int32   `json:"stop_code"`
	StopID          string  `json:"stop_id"`
	StopDescr       string  `json:"stop_descr"`
	StopDescrEng    string  `json:"stop_descr_eng"`
	StopStreet      string  `json:"stop_street"`
	StopStreetEng   string  `json:"stop_street_eng"`
	StopHeading     int32   `json:"stop_heading"`
	StopLat         float64 `json:"stop_lat"`
	StopLng         float64 `json:"stop_lng"`
	StopType        int8    `json:"stop_type"`
	StopAmea        int8    `json:"stop_amea"`
	Destinations    string  `json:"destinations"`
	DestinationsEng string  `json:"destinations_eng"`
	StopSenu        int16   `json:"senu"`
}

type StopDtoM struct {
	StopSenu  int16  `json:"stopsenu"`
	StopCode  int32  `json:"stopcode"`
	StopID    string `json:"stopid"`
	StopDescr string `json:"stopdescr"`
}
