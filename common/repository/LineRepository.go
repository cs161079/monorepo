package repository

import (
	"errors"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type lineRepository struct {
	DB *gorm.DB
}

type LineRepository interface {
	SelectByCode(lineCode int32) (*models.Line, error)
	Insert(line *models.Line) (*models.Line, error)
	InsertArray([]models.Line) ([]models.Line, error)
	Update(line *models.Line) (*models.Line, error)
	LineList01() ([]models.LineDto01, error)
	DeleteAll() error
	WithTx(*gorm.DB) lineRepository
}

func NewLineRepository(connection *gorm.DB) LineRepository {
	return lineRepository{
		DB: connection,
	}
}

// withTx creates a new repository instance with the given transaction
func (r lineRepository) WithTx(tx *gorm.DB) lineRepository {
	if tx == nil {
		logger.WARN("Database Tranction not exist.")
		return r
	}
	r.DB = tx
	return r
}

func (r lineRepository) SelectByCode(lineCode int32) (*models.Line, error) {
	var result models.Line
	dbResults := r.DB.Table(db.LINETABLE).Where("line_code = ?", lineCode).Find(&result)
	if dbResults.RowsAffected == 0 {
		dbResults.Error = gorm.ErrRecordNotFound
	}
	if dbResults.Error != nil {
		if errors.Is(dbResults.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else if errors.Is(dbResults.Error, MyError{}) {
			//panic(fmt.Sprintln("Database Error ", results.Error.Error()))
			return nil, dbResults.Error
		}
	}
	return &result, nil
}

func (r lineRepository) Insert(line *models.Line) (*models.Line, error) {
	trxRes := r.DB.Table(db.LINETABLE).Create(line)
	if trxRes.Error != nil {
		return nil, trxRes.Error
	}
	return line, nil
}

func (r lineRepository) Update(line *models.Line) (*models.Line, error) {
	trxRes := r.DB.Table(db.LINETABLE).Save(line)
	if trxRes.Error != nil {
		return nil, trxRes.Error
	}
	return line, nil
}

func (r lineRepository) LineList01() ([]models.LineDto01, error) {
	var result []models.LineDto01
	res := r.DB.Table(db.LINETABLE).Where("mld_master=?", 1).Order("line_id").Find(&result)
	if res != nil {
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return result, nil
}

func (r lineRepository) DeleteAll() error {
	if err := r.DB.Table(db.LINETABLE).Where("1=1").Delete(&models.Line{}).Error; err != nil {
		return err
	}
	return nil
}

func (r lineRepository) InsertArray(entityArr []models.Line) ([]models.Line, error) {
	if err := r.DB.Table(db.LINETABLE).Save(entityArr).Error; err != nil {
		return nil, err
	}
	return entityArr, nil
}

type MyError struct {
}

func (t MyError) Error() string {
	return "This is my critical Error"
}
