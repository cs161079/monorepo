package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"gorm.io/gorm"
)

type RouteRepository interface {
	SelectByCode(int32) (*models.Route, error)
	SelectByLineCode(int32) (*[]models.Route, error)
	Insert(models.Route) (*models.Route, error)
	InsertArray([]models.Route) ([]models.Route, error)
	Update(models.Route) (*models.Route, error)
	List01() ([]models.Route, error)
	WithTx(*gorm.DB) routeRepository
	DeleteAll() error
}

type routeRepository struct {
	DB *gorm.DB
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

func (r routeRepository) SelectByLineCode(lineCode int32) (*[]models.Route, error) {
	var selectedVal []models.Route
	res := r.DB.Table(db.ROUTETABLE).Where("line_code = ?", lineCode).Find(&selectedVal)
	if res.Error != nil {
		return nil, res.Error
	}
	return &selectedVal, nil
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

func NewRouteRepository(dbConnection *gorm.DB) RouteRepository {
	return routeRepository{
		DB: dbConnection,
	}
}
