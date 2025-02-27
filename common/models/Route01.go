package models

// ********* Struct for Route Details **************
// ************** Database Entity ********************
type Route01 struct {
	Id           int64   `json:"id" gorm:"primaryKey"`
	Rt_code      int32   `json:"route_code" gorm:"index:route01_code_un,unique"`
	Routed_x     float32 `json:"routed_x"`
	Routed_y     float32 `json:"routed_y"`
	Routed_order int16   `json:"routed_order" gorm:"index:route01_code_un,unique"`
}

func (Route01) TableName() string {
	return "Route01"
}

//**********************************************************

// ********* Struct for Route Details OASA **************
type Route01Oasa struct {
	Routed_x     float32 `oasa:"routed_x"`
	Routed_y     float32 `oasa:"routed_y"`
	Routed_order int16   `oasa:"routed_order"`
}
