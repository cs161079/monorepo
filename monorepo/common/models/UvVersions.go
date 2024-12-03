package models

const (
	UVERSIONS_LINE       = "LINES"
	UVERSIONS_ROUTE      = "ROUTES"
	UVERSIONS_STOP       = "STOPS"
	UVERSIONS_ROUTESTOPS = "ROUTE STOPS"
	UVERSIONS_SCHED_CATS = "SCHED_CATS"
	SCHED_ENTRIES        = "SCHED_ENTRIES"
)

type UVersionsOasa struct {
	Uv_descr          string `json:"uv_descr" oasa:"UV_DESCR"`
	Uv_lastupdatelong int64  `json:"uv_lastupdatelong" oasa:"UV_LASTUPDATELONG"`
}

type UVersions struct {
	Uv_descr          string `json:"uv_descr" gorm:"primaryKey"`
	Uv_lastupdatelong int64  `json:"uv_lastupdatelong"`
}
