package service

import (
	"encoding/json"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type RouteService interface {
	WithTrx(*gorm.DB) routeService
	DeleteAll() error
	DeleteRoute01() error
	DeleteRoute02() error
	InsertArray([]models.Route) ([]models.Route, error)
	InserChunkArray(chunSize int, allData []models.Route) error
	Route02InsertArr([]models.Route02) ([]models.Route02, error)
	Route02InsertChunkArray(chunkSize int, allData []models.Route02) error
	Route01InsertArr([]models.Route01) ([]models.Route01, error)
	Route01InsertChunkArray(chunkSize int, allData []models.Route01) error

	SelectFirstRouteByLinecodeWithStops(line_code int32) (*models.RouteDto, error)
	SelectRouteWithStops(int32) (*models.RouteDto, error)
	SelectRouteDetails(int32) ([]models.Route01, error)
	SelectRouteStop(int32) ([]models.Route02Dto, error)
	PassengersCount(int, int) (*models.BusCapacityDt02, error)

	// ---------------- For Trip Planner -----------------
	RouteList() ([]models.RouteWithLine, error)
	RouteStopList(int32) ([]models.Route02, error)
	RouteSelect(int32) (*models.Route, error)
}

type routeService struct {
	repo        repository.RouteRepository
	repo02      repository.Route02Repository
	repo01      repository.Route01Repository
	mapper01    mapper.Route01Mapper
	routeMapper mapper.RouteMapper
}

func NewRouteService(repo repository.RouteRepository,
	repo01 repository.Route01Repository,
	repo02 repository.Route02Repository) RouteService {
	return &routeService{
		repo:        repo,
		repo02:      repo02,
		repo01:      repo01,
		mapper01:    mapper.NewRouteDetailMapper(),
		routeMapper: mapper.NewRouteMapper(),
	}
}

func (s routeService) WithTrx(trxHandle *gorm.DB) routeService {
	s.repo = s.repo.WithTx(trxHandle)
	return s
}

func (s routeService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s routeService) DeleteRoute02() error {
	return s.repo02.DeleteRoute02()
}

func (s routeService) Route02InsertArr(entityArr []models.Route02) ([]models.Route02, error) {
	err := s.repo02.InsertRoute02Arr(entityArr)
	if err != nil {
		return nil, err
	}
	return entityArr, nil
}

func (s routeService) Route02InsertChunkArray(chunkSize int, allData []models.Route02) error {
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}
	for {
		_, err := s.Route02InsertArr(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι διαδρομές από %d έως %d.", stratIndex, endIndex-1))
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.Route02InsertArr(allData[stratIndex:])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (s routeService) InsertArray(entityArr []models.Route) ([]models.Route, error) {
	return s.repo.InsertArray(entityArr)
}

func (s routeService) InserChunkArray(chunkSize int, allData []models.Route) error {
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}
	for {
		_, err := s.InsertArray(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι διαδρομές από %d έως %d.", stratIndex, endIndex-1))
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.InsertArray(allData[stratIndex:])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (s routeService) Route01InsertChunkArray(chunkSize int, allData []models.Route01) error {
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}

	for {
		_, err := s.Route01InsertArr(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
		//logger.INFO(fmt.Sprintf("Προστέθηκαν οι λεπτομερειες διαδρομών από %d έως %d.", stratIndex, endIndex-1))
		stratIndex = endIndex
		endIndex = stratIndex + chunkSize
		if stratIndex > len(allData)-1 {
			//logger.INFO("Η εισαγωγή λεπτομερειών διαδρομών ολοκληρώθηκε.")
			break
		} else if endIndex > len(allData)-1 {
			_, err := s.Route01InsertArr(allData[stratIndex:])
			if err != nil {
				//logger.ERROR(fmt.Sprintf("Σφάλμα κατά την προσθήκη λεπτομερειών διαδρομών από %d έως τέλος.", stratIndex))
				//txt.Rollback()
				return err
			}
			break
		}
	}
	return nil
}

func (s routeService) Route01InsertArr(entityArr []models.Route01) ([]models.Route01, error) {
	return s.repo01.InsertRoute01Arr(entityArr)
}

func (s routeService) DeleteRoute01() error {
	return s.repo01.Delete()
}

func (s routeService) SelectFirstRouteByLinecodeWithStops(line_code int32) (*models.RouteDto, error) {
	origData, err := s.repo.SelectByLineCodeWithStops(line_code)
	if err != nil {
		return nil, err
	}
	return s.routeMapper.RouteToRouteDto(*origData)
}

func (s routeService) SelectRouteWithStops(routeCode int32) (*models.RouteDto, error) {

	// Get Data from Database
	origData, err := s.repo.SelectByRouteCodeWithStops(routeCode)
	if err != nil {
		return nil, err
	}
	return s.routeMapper.RouteToRouteDto(*origData)
}

func (s routeService) SelectRouteDetails(routeCode int32) ([]models.Route01, error) {
	return s.repo01.SelectByCode(routeCode)
}

func (s routeService) SelectRouteStop(routecode int32) ([]models.Route02Dto, error) {
	return s.repo02.SelectRouteStops(routecode)
}

func (s routeService) PassengersCount(busID int, routeID int) (*models.BusCapacityDt02, error) {
	dbRec, err := s.repo.PassengersCount(busID, routeID)
	if err != nil {
		return nil, err
	}
	var result models.BusCapacityDt02 = models.BusCapacityDt02{}
	bts, err := json.Marshal(dbRec)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bts, &result)
	return &result, nil
}

func (s routeService) RouteList() ([]models.RouteWithLine, error) {
	return s.repo.RouteList()
}

func (s routeService) RouteStopList(routeCode int32) ([]models.Route02, error) {
	return s.repo.RouteStopList(routeCode)
}

func (s routeService) RouteSelect(routeCode int32) (*models.Route, error) {
	return s.repo.SelectByCode(routeCode)
}
