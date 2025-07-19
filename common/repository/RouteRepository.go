package repository

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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
	PassengersCount(busID int, routeID int) (*models.Bus_Capacity, error)

	// ---------------------- For Trip Planner -----------------------------------
	RouteList() ([]models.RouteWithLine, error)
	RouteStopList(int32) ([]models.Route02, error)
	//----------------------------------------------------------------------------
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
	dbResult := r.DB.Select("route.route_code", "line.line_code", "route.route_descr", "line.line_id, line.line_type").Table(db.ROUTESTOPSTABLE).Joins(
		"LEFT JOIN route on route02.rt_code=route.route_code").Joins(
		"LEFT JOIN line on route.ln_code=line.line_code").Where(
		"route02.stp_code=?", stop_code).Find(&result)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return result, nil
}

func (r routeRepository) PassengersCount(busID int, routeID int) (*models.Bus_Capacity, error) {
	now := time.Now()
	from := now.Add(-1 * 5 * time.Minute)
	to := now.Add(1 * 5 * time.Minute)
	strFrom := from.Format("2006-01-02 15:04:05") // YYYY-MM-DD HH:mm:ss format
	strTo := to.Format("2006-01-02 15:04:05")

	var result models.Bus_Capacity
	dbResult := r.DB.Table(db.BUSCAPACITY).Where(
		"date_time BETWEEN ? and ? AND bus_id = ? AND route_id = ?", strFrom, strTo, busID, routeID).Order("date_time desc").Limit(1).Find(&result)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &result, nil
}

func (r routeRepository) RouteList() ([]models.RouteWithLine, error) {
	var dbData []models.RouteWithLine = make([]models.RouteWithLine, 0)
	dbResults := r.DB.Select("route.route_code, route.ln_code, line.line_id, route.route_type, route.route_descr").
		Table(db.ROUTETABLE).Joins("LEFT JOIN line on route.ln_code=line.line_code").
		Order("route.route_code, route.ln_code").
		Find(&dbData)
	if dbResults.Error != nil {
		return nil, dbResults.Error
	}
	return dbData, nil
}

func (r routeRepository) RouteStopList(routeCode int32) ([]models.Route02, error) {
	var dbData []models.Route02 = make([]models.Route02, 0)
	dbResult := r.DB.Table(db.ROUTESTOPSTABLE).
		Where("rt_code=?", routeCode).
		Order("senu").Find(&dbData)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return dbData, nil
}
