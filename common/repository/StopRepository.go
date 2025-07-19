package repository

import (
	"fmt"
	"math"
	"sort"

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
	SelectClosestStops(float64, float64, float32, float32) ([]models.StopDto, error)
	SelectClosestStops02(float64, float64) ([]models.StopDto, error)
	WithTx(*gorm.DB) stopRepository

	SelectAll() ([]models.Stop, error)
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

func (r stopRepository) SelectClosestStops(lat float64, long float64, from float32, to float32) ([]models.StopDto, error) {
	var resultList []models.StopDto
	var subQuery = r.DB.Table("stop s").Select("stop_code, stop_descr, stop_street," +
		fmt.Sprintf("round(haversine_distance(%f, %f, s.stop_lat, s.stop_lng), 2)", lat, long) +
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

// Haversine formula to calculate the great-circle distance between two points
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (r stopRepository) SelectClosestStops02(latitude float64, longtitude float64) ([]models.StopDto, error) {
	sqlDb, err := r.DB.DB()
	if err != nil {
		return nil, err
	}
	rows, err := sqlDb.Query("SELECT s.stop_code, s.stop_descr, s.stop_street, s.stop_lat, s.stop_lng FROM stop s")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stops []models.StopDto = make([]models.StopDto, 0)

	for rows.Next() {
		var stopDto models.StopDto
		// var stop models.Stop
		if err := rows.Scan(&stopDto.Stop_code, &stopDto.Stop_descr, &stopDto.Stop_street, &stopDto.Stop_lat, &stopDto.Stop_lng); err != nil {
			continue
		}
		// mapper.MapStruct(stop, &stopDto)
		stopDto.Distance = haversine(latitude, longtitude, stopDto.Stop_lat, stopDto.Stop_lng)
		if stopDto.Distance < 1.5 {
			stops = append(stops, stopDto)
		}
	}

	// Sort by distance
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].Distance < stops[j].Distance
	})

	// Return the closest 10 stops
	if len(stops) > 20 {
		stops = stops[:20]
	}
	return stops, nil
}

func (s stopRepository) SelectAll() ([]models.Stop, error) {
	var dbData []models.Stop = make([]models.Stop, 0)
	dbResults := s.DB.Table(db.STOPTABLE).Find(&dbData)
	if dbResults.Error != nil {
		return nil, dbResults.Error
	}
	return dbData, nil
}
