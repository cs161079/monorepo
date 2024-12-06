package repository

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type Route02Repository interface {
	SelectByCode(int32, int64, int16) (*models.Route02, error)
	DeleteStopByRoute(int32) error
	InsertRoute02(models.Route02) error
	InsertRoute02Arr([]models.Route02) error
	UpdateRoute02(models.Route02) error
	DeleteRoute02() error
}

type route02Repository struct {
	DB *gorm.DB
}

func (r route02Repository) SelectByCode(routecode int32, stopcode int64, senu int16) (*models.Route02, error) {
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

func NewRoute02Repository(dbConnection *gorm.DB) Route02Repository {
	return route02Repository{
		DB: dbConnection,
	}
}
