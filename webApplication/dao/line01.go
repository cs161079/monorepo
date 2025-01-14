package dao

type Line01 struct {
	Line_Code  int32  `json:"line_code"`
	Line_Descr string `json:"line_descr"`
	Routes     []Route01
}

type Route01 struct {
	Route_code  int32
	Route_Descr string
}

// type Route02 struct {
// 	Route_Code int32
// }

type Stop01 struct {
	Stop_code  int32
	Stop_descr string
}
