package dao

type Line struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	Line_Code  int32     `json:"line_code"`
	Line_Descr string    `json:"line_descr"`
	Routes     []Route01 `json:"routes" gorm:"foreignKey:line_code;reference:line_code"`
}

func (Line) TableName() string {
	return "line"
}

type Route01 struct {
	Route_code  int32  `json:"route_code"`
	Route_Descr string `json:"route_descr"`
	Line_Code   int32  `json:"line_code"`
}

func (Route01) TableName() string {
	return "route"
}

// type Route02 struct {
// 	Route_Code int32
// }

type Stop01 struct {
	Stop_code  int32
	Stop_descr string
}
