package repository

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type StopRepository interface {
	SelectByCode(int32) (*models.Stop, error)
	Insert(models.Stop) (*models.Stop, error)
	InsertArray([]models.Stop) ([]models.Stop, error)
	Update(models.Stop) (*models.Stop, error)
	List01(int32) (*[]models.Stop, error)
	DeleteAll() error
	SelectClosestStops(models.Point, float32, float32) ([]models.StopDto, error)
	WithTx(*gorm.DB) stopRepository
}

type stopRepository struct {
	DB *gorm.DB
}

func NewStopRepository(connection *gorm.DB) StopRepository {
	return stopRepository{
		DB: connection,
	}
}

func (r stopRepository) WithTx(tx *gorm.DB) stopRepository {
	if tx == nil {
		logger.WARN("Database Tranction not exist.")
		return r
	}
	r.DB = tx
	return r
}

func (r stopRepository) SelectByCode(stopCode int32) (*models.Stop, error) {
	var selectedVal models.Stop
	res := r.DB.Table(db.STOPTABLE).Where("stop_code = ?", stopCode).Find(&selectedVal)
	if res.Error != nil {
		fmt.Println(res.Error.Error())
		return nil, res.Error
	}
	return &selectedVal, nil
}

func (r stopRepository) Insert(busStop models.Stop) (*models.Stop, error) {
	res := r.DB.Table(db.STOPTABLE).Create(&busStop)
	if res.Error != nil {
		return nil, res.Error
	}
	return &busStop, nil
}

func (r stopRepository) Update(busStop models.Stop) (*models.Stop, error) {
	res := r.DB.Table(db.STOPTABLE).Save(&busStop)
	if res.Error != nil {
		return nil, res.Error
	}
	return &busStop, nil
}

func (r stopRepository) List01(routeCode int32) (*[]models.Stop, error) {
	var result []models.Stop
	res := r.DB.Table(db.STOPTABLE).
		Select("stop.*, "+
			"routestops.senu").
		Joins("LEFT JOIN routestops ON stop.stop_code=routestops.stop_code").
		Where("routestops.route_code=?", routeCode).Order("senu").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return &result, nil
}

func (r stopRepository) DeleteAll() error {
	if err := r.DB.Table(db.STOPTABLE).Where("1=1").Delete(&models.Stop{}).Error; err != nil {
		return err
	}
	return nil
}

func (r stopRepository) SelectClosestStops(point models.Point, from float32, to float32) ([]models.StopDto, error) {
	var resultList []models.StopDto
	var subQuery = r.DB.Table("stop s").Select("stop_code, stop_descr, stop_street," +
		fmt.Sprintf("round(haversine_distance(%f, %f, s.stop_lat, s.stop_lng), 2)", point.Lat, point.Long) +
		" AS distance")

	if err := r.DB.Table("(?) as b", subQuery).Select(" b. stop_code, b.stop_descr, b.stop_street, b.distance").
		Where(
			fmt.Sprintf(
				"distance > %f AND distance <= %f", from, to)).
		Order("distance").
		Find(&resultList).Error; err != nil {
		return nil, err
	}
	return resultList, nil

}

func (r stopRepository) InsertArray(entityArray []models.Stop) ([]models.Stop, error) {
	if err := r.DB.Table(db.STOPTABLE).Save(entityArray).Error; err != nil {
		return nil, err
	}
	return entityArray, nil
}
