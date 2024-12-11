package repository

import (
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/models"

	"gorm.io/gorm"
)

type Route01Repository interface {
	InsertRoute01Arr([]models.Route01) ([]models.Route01, error)
	Delete() error
}

type route01Repository struct {
	DB *gorm.DB
}

func NewRoute01Repository(connection *gorm.DB) Route01Repository {
	return route01Repository{
		DB: connection,
	}
}

func (r route01Repository) InsertRoute01Arr(entityArr []models.Route01) ([]models.Route01, error) {
	res := r.DB.Table(db.ROUTEDETAILTABLE).Save(entityArr)
	if res.Error != nil {
		return nil, res.Error
	}
	return entityArr, nil
}

func (r route01Repository) Delete() error {
	if err := r.DB.Table(db.ROUTEDETAILTABLE).Where("1=1").Delete(&models.Route01{}).Error; err != nil {
		//trans.Rollback()
		return err
	}
	return nil
}
