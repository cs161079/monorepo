package models

import "github.com/cs161079/monorepo/common/db"

// ********* Struct for Route Details **************
// ************** Database Entity ********************
// type Route01 struct {
// 	Id           int64   `json:"id" gorm:"primaryKey"`
// 	Rt_code      int32   `json:"route_code" gorm:"index:route01_code_un,unique"`
// 	Routed_x     float32 `json:"routed_x"`
// 	Routed_y     float32 `json:"routed_y"`
// 	Routed_order int16   `json:"routed_order" gorm:"index:route01_code_un,unique"`
// }

type Route01 struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	RtCode      int32   `json:"route_code" gorm:"column:rt_code"` // FK to route.route_code
	RoutedX     float64 `json:"routed_x" gorm:"column:routed_x"`
	RoutedY     float64 `json:"routed_y" gorm:"column:routed_y"`
	RoutedOrder int16   `json:"routed_order" gorm:"column:routed_order"`

	Route Route `json:"route" gorm:"foreignKey:RtCode;references:RouteCode"`
}

func (Route01) TableName() string {
	return db.ROUTEDETAILTABLE
}

//**********************************************************

// ********* Struct for Route Details OASA **************
type Route01Oasa struct {
	Routed_x     float32 `oasa:"routed_x"`
	Routed_y     float32 `oasa:"routed_y"`
	Routed_order int16   `oasa:"routed_order"`
}
