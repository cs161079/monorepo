package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type RouteRepository interface {
	SelectByCode(int32) (*models.Route, error)
	SelectByLineCodeWithStops(int32) (*models.Route, error)
	SelectByRouteCodeWithStops(int32) (*models.Route, error)
	Insert(models.Route) (*models.Route, error)
	InsertArray([]models.Route) ([]models.Route, error)
	Update(models.Route) (*models.Route, error)
	// Returns LINE IDs and ROUTES that pass through this stop.
	//
	// @input stop_code int32 STOP code.
	//
	// @returns []models.StopArrival, error.
	ExtraArrivalInfo(int32) ([]models.StopArrival, error)
	List01() ([]models.Route, error)
	WithTx(*gorm.DB) routeRepository
	DeleteAll() error
}

type routeRepository struct {
	DB *gorm.DB
}

func NewRouteRepository(connection *gorm.DB) RouteRepository {
	return routeRepository{
		DB: connection,
	}
}

func (r routeRepository) SelectByCode(routeCode int32) (*models.Route, error) {
	var selectedVal models.Route
	dbRes := r.DB.Table(db.ROUTETABLE).Where("route_code = ?", routeCode).Find(&selectedVal)
	if dbRes != nil {
		if dbRes.Error != nil {
			// fmt.Println(r.Error.Error())
			return nil, dbRes.Error
		}
		if dbRes.RowsAffected == 0 {
			//logger.WARN(fmt.Sprintf("BUS ROUTE NOT FOUND [ROUTE_CODE: %d]", routeCode))
			return nil, nil
		}
	}
	return &selectedVal, nil
}

func (r routeRepository) SelectByLineCodeWithStops(lineCode int32) (*models.Route, error) {
	var result models.Route
	dbResults := r.DB.Preload("Route02s", func(db *gorm.DB) *gorm.DB {
		return db.Order("route02.senu")
	}).Preload("Route02s.Stop").Where("ln_Code = ?", lineCode).Order("route_code").First(&result)
	if dbResults.Error != nil {
		if errors.Is(dbResults.Error, gorm.ErrRecordNotFound) {
			return nil, models.NewError(dbResults.Error.Error(),
				fmt.Sprintf("No Route found with code %d.", lineCode), http.StatusNotFound)
		}
		return nil, dbResults.Error
	}

	return &result, nil
}

func (r routeRepository) SelectByRouteCodeWithStops(routeCd int32) (*models.Route, error) {
	var result models.Route
	dbResults := r.DB.Preload("Route02s", func(db *gorm.DB) *gorm.DB {
		return db.Order("route02.senu")
	}).Preload("Route02s.Stop").Where("route.route_code = ?", routeCd).First(&result)
	if dbResults.Error != nil {
		if errors.Is(dbResults.Error, gorm.ErrRecordNotFound) {
			return nil, models.NewError(dbResults.Error.Error(),
				fmt.Sprintf("No Route found with code %d.", routeCd), http.StatusNotFound)
		}
		return nil, dbResults.Error
	}

	return &result, nil
}

func (r routeRepository) Insert(input models.Route) (*models.Route, error) {
	res := r.DB.Table(db.ROUTETABLE).Create(&input)
	if res.Error != nil {
		return nil, res.Error
	}
	return &input, nil
}

func (r routeRepository) Update(input models.Route) (*models.Route, error) {
	res := r.DB.Table(db.ROUTETABLE).Create(&input)
	if res.Error != nil {
		return nil, res.Error
	}
	return &input, nil
}

func (r routeRepository) List01() ([]models.Route, error) {
	var result []models.Route
	res := r.DB.Table(db.ROUTETABLE).Order("route_code").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (r routeRepository) DeleteAll() error {
	if err := r.DB.Table(db.ROUTETABLE).Where("1=1").Delete(&models.Route{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}

func (r routeRepository) WithTx(tx *gorm.DB) routeRepository {
	if tx == nil {
		logger.WARN("Database Tranction not exist.")
		return r
	}
	r.DB = tx
	return r
}

func (r routeRepository) InsertArray(entiryArr []models.Route) ([]models.Route, error) {
	if err := r.DB.Table(db.ROUTETABLE).Save(entiryArr).Error; err != nil {
		return nil, err
	}
	return entiryArr, nil
}

func (r routeRepository) ExtraArrivalInfo(stop_code int32) ([]models.StopArrival, error) {
	var result []models.StopArrival
	dbResult := r.DB.Select("route.route_code", "line.line_descr", "line.line_id").Table(db.ROUTESTOPSTABLE).Joins(
		"LEFT JOIN route on route02.rt_code=route.route_code").Joins(
		"LEFT JOIN line on route.ln_code=line.line_code").Where(
		"route02.stp_code=?", stop_code).Find(&result)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return result, nil
}
