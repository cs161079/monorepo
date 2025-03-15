package repository

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type Route02Repository interface {
	SelectByCode(int32, int32, int16) (*models.Route02, error)
	SelectRouteStops(int32) ([]models.Route02Dto, error)
	DeleteStopByRoute(int32) error
	InsertRoute02(models.Route02) error
	InsertRoute02Arr([]models.Route02) error
	UpdateRoute02(models.Route02) error
	DeleteRoute02() error
}

type route02Repository struct {
	DB *gorm.DB
}

func NewRoute02Repository(connection *gorm.DB) Route02Repository {
	return route02Repository{
		DB: connection,
	}
}

func (r route02Repository) SelectByCode(routecode int32, stopcode int32, senu int16) (*models.Route02, error) {
	var result models.Route02
	dbRes := r.DB.Table(db.ROUTESTOPSTABLE).
		Where("route_code=? and stop_code=? and senu=?", routecode, stopcode, senu).Find(&result)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	return &result, nil
}

func (r route02Repository) DeleteStopByRoute(routeCode int32) error {
	var routeStops []models.Route02
	res := r.DB.Table(db.ROUTESTOPSTABLE).Where("route_code=?", routeCode).Delete(&routeStops)
	if res.Error != nil {
		// logger.ERROR(r.Error.Error())
		return res.Error
	}
	logger.INFO(fmt.Sprintf("DELETED ROWS %d", res.RowsAffected))
	return nil
}

func (r route02Repository) InsertRoute02(input models.Route02) error {
	res := r.DB.Table(db.ROUTESTOPSTABLE).Create(&input)
	if res.Error != nil {
		return res.Error
	}
	//logger.INFO(fmt.Sprintf("STOP [%d] SAVED SUCCESFULL IN ROUTE [%d].", input.Stop_code, input.Route_code))
	return nil
}

func (r route02Repository) UpdateRoute02(input models.Route02) error {
	res := r.DB.Table(db.ROUTESTOPSTABLE).Create(&input)
	if res.Error != nil {
		return res.Error
	}
	//logger.INFO(fmt.Sprintf("STOP [%d] SAVED SUCCESFULL IN ROUTE [%d].", input.Stop_code, input.Route_code))
	return nil
}

func (r route02Repository) DeleteRoute02() error {
	if err := r.DB.Table(db.ROUTESTOPSTABLE).Where("1=1").Delete(&models.Route02{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r route02Repository) InsertRoute02Arr(entityArr []models.Route02) error {
	res := r.DB.Table(db.ROUTESTOPSTABLE).Save(entityArr)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r route02Repository) SelectRouteStops(routeCode int32) ([]models.Route02Dto, error) {
	var result []models.Route02Dto
	dbResult := r.DB.Select("stop.stop_code, stop.stop_descr, stop.stop_lat, stop.stop_lng, route02.senu").Table(db.ROUTESTOPSTABLE).Joins("LEFT JOIN stop on route02.stp_code=stop.stop_code").Where("route02.rt_code=?", routeCode).Order("route02.senu").Find(&result)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return result, nil
}
