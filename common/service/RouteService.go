package service

import (
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type RouteService interface {
	Post(models.Route) (*models.Route, error)
	DeleteAll() error
	WithTrx(*gorm.DB) routeService
	PostRoute02(models.Route02) (*models.Route02, error)
	InsertArray([]models.Route) ([]models.Route, error)
	Route02Insert(models.Route02) (*models.Route02, error)
	Route02InsertArr([]models.Route02) ([]models.Route02, error)
	Route01InsertArr([]models.Route01) ([]models.Route01, error)
	DeleteRoute01() error
	DeleteRoute02() error
	List01() ([]models.Route, error)
	GetMapper01() mapper.Route01Mapper
}

type routeService struct {
	repo     repository.RouteRepository
	repo02   repository.Route02Repository
	repo01   repository.Route01Repository
	mapper01 mapper.Route01Mapper
}

func NewRouteService(repo repository.RouteRepository,
	repo01 repository.Route01Repository,
	repo02 repository.Route02Repository, route01Mapper mapper.Route01Mapper) RouteService {
	return routeService{
		repo:     repo,
		repo02:   repo02,
		repo01:   repo01,
		mapper01: route01Mapper,
	}
}

func (s routeService) WithTrx(trxHandle *gorm.DB) routeService {
	s.repo = s.repo.WithTx(trxHandle)
	return s
}

func (s routeService) Post(entity models.Route) (*models.Route, error) {
	routeFound, err := s.repo.SelectByCode(entity.Route_Code)
	if err != nil {
		return nil, err
	}
	if routeFound != nil {
		return s.repo.Update(entity)
	} else {
		return s.repo.Insert(entity)
	}
}

func (s routeService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s routeService) PostRoute02(entity models.Route02) (*models.Route02, error) {
	foundRoute02, err := s.repo02.SelectByCode(entity.Route_code, entity.Stop_code, entity.Senu)
	if err != nil {
		return nil, err
	}
	if foundRoute02 == nil {
		err = s.repo02.InsertRoute02(entity)
	} else {
		err = s.repo02.UpdateRoute02(entity)
	}
	if err != nil {
		return nil, err
	}
	return foundRoute02, nil
}

func (s routeService) DeleteRoute02() error {
	return s.repo02.DeleteRoute02()
}

func (s routeService) Route02Insert(entity models.Route02) (*models.Route02, error) {
	err := s.repo02.InsertRoute02(entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s routeService) Route02InsertArr(entityArr []models.Route02) ([]models.Route02, error) {
	err := s.repo02.InsertRoute02Arr(entityArr)
	if err != nil {
		return nil, err
	}
	return entityArr, nil
}

func (s routeService) InsertArray(entityArr []models.Route) ([]models.Route, error) {
	return s.repo.InsertArray(entityArr)
}

func (s routeService) Route01InsertArr(entityArr []models.Route01) ([]models.Route01, error) {
	return s.repo01.InsertRoute01Arr(entityArr)
}

func (s routeService) List01() ([]models.Route, error) {
	return s.repo.List01()
}

func (s routeService) GetMapper01() mapper.Route01Mapper {
	return s.mapper01
}

func (s routeService) DeleteRoute01() error {
	return s.repo01.Delete()
}
