package service

import (
	"github.com/cs161079/monorepo/common/mapper"
	models "github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
)

type UVersionService interface {
	Select(string) (*models.UVersions, error)
	Post(*models.UVersions) error
}

type uVersionsService struct {
	Repo   repository.UVersionRepository
	Rest   RestService
	Mapper mapper.UversionMapper
}

func NewuVersionService(repo repository.UVersionRepository) UVersionService {
	return uVersionsService{
		Repo: repo,
	}
}

func (s uVersionsService) GetUversionWeb() {
	response := s.Rest.OasaRequestApi00("getUversions", nil)
	if response.Error != nil {
		// Εδώ προκείπτει Error από το Request
		// Κάτι πρέπει να κάνουμε.
	}
	// var arrVersion []models.UVersions = make([]models.UVersions, 0)
	// for index, int := range response.Data.([]interface{}) {
	// 	uVersion := s.Mapper.OasaToUVersions(s.Mapper.GeneralUVersions(response.Data))

	// }

}

func (s uVersionsService) Post(entity *models.UVersions) error {
	var dbEntity *models.UVersions = nil
	var err error = nil
	dbEntity, err = s.Repo.Select(entity.Uv_descr)
	if err != nil {
		return err
	}
	isNew := dbEntity == nil
	if isNew {
		return s.Repo.Create(entity)
	} else {
		return s.Repo.Update(entity)
	}
}

func (s uVersionsService) Select(uVersion string) (*models.UVersions, error) {
	return s.Repo.Select(uVersion)
}
