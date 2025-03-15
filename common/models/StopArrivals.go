package models

type StopArrivalOasa struct {
	Btime2     int8  `json:"time" oasa:"btime2"`
	Route_code int32 `json:"route_code" oasa:"route_code"`
	Veh_code   int32 `json:"veh_code" oasa:"veh_code"`
}

type StopArrival struct {
	Btime2     int8   `json:"time"`
	Route_code int32  `json:"route_code"`
	Line_descr string `json:"line_descr"`
	Veh_code   int32  `json:"veh_code"`
	Line_id    string `json:"line_id"`
	NextTime   int8   `json:"next_time"`
	Last_Time  int8   `json:"last_time"`
}

type BusLocation struct {
	Veh_code   int32   `json:"veh_code" oasa:"VEH_NO"`
	Cs_date    string  `json:"cs_date" oasa:"CS_DATE"`
	Cs_lat     float64 `json:"cs_lat" oasa:"CS_LAT"`
	Cs_lng     float64 `json:"cs_lng" oasa:"CS_LNG"`
	Route_code int32   `json:"route_code" oasa:"ROUTE_CODE"`
}
