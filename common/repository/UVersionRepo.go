package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"

	"gorm.io/gorm"
)

type UVersionRepository interface {
	Create(entity *models.UVersions) error
	Update(entity *models.UVersions) error
	SelectAll() ([]models.UVersions, error)
	Select(string) (*models.UVersions, error)
	WithTx(*gorm.DB) uVersionRepository
}

type uVersionRepository struct {
	DB *gorm.DB
}

// // withTx creates a new repository instance with the given transaction
func (r uVersionRepository) WithTx(tx *gorm.DB) uVersionRepository {
	if tx == nil {
		logger.WARN("Database transaction does not exist.")
		return r
	}
	r.DB = tx
	return r
}

// // Update modifies an existing entity in the database
func (r uVersionRepository) Update(entity *models.UVersions) error {
	return r.DB.Table(db.SYNCVERSIONSTABLE).Save(entity).Error
}

func (r uVersionRepository) SelectAll() ([]models.UVersions, error) {
	var res []models.UVersions

	if err := r.DB.Table(db.SYNCVERSIONSTABLE).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil

}

// // Create adds a new entity to the database
func (r uVersionRepository) Create(entity *models.UVersions) error {
	return r.DB.Create(entity).Error
}

func (r uVersionRepository) Select(uVersion string) (*models.UVersions, error) {
	var resultRec models.UVersions = models.UVersions{}
	dbRes := r.DB.Table("syncversions").Where("uv_descr=?", uVersion).Find(&resultRec)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}
	return &resultRec, nil
}

func NewUversionRepository(db *gorm.DB) UVersionRepository {
	return uVersionRepository{DB: db}
}
