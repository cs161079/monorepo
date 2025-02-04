package service

import (
	"context"
	"fmt"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/utils"

	"gorm.io/gorm"
)

type LineService interface {
	WithTrx(*gorm.DB) lineService
	InsertArray([]models.Line) ([]models.Line, error)
	InsertChunkArray(chunkSize int, allData []models.Line) error
	DeleteAll() error
	GetLineList() ([]models.LineDto01, error)
	SelectByLineCode(lineCode string) (*models.LineDto, error)

	InsertLine(line *models.Line) (*models.Line, error)
	PostLine(line *models.Line) (*models.Line, error)
	PostLineArray(context.Context, []models.Line) ([]models.Line, error)

	GetMapper() mapper.LineMapper
}
type lineService struct {
	repo   repository.LineRepository
	mapper mapper.LineMapper
}

func NewLineService(repo repository.LineRepository) LineService {
	return lineService{
		repo:   repo,
		mapper: mapper.NewLineMapper(),
	}
}

func (s lineService) GetMapper() mapper.LineMapper {
	return s.mapper
}

func (s lineService) SelectByLineCode(line_code string) (*models.LineDto, error) {
	num32, err := utils.StrToInt32(line_code)
	if err != nil {
		return nil, fmt.Errorf("Error on converting String to numbe. %s", err.Error())
	}

	line, err := s.repo.SelectByCode(*num32)
	if err != nil {
		return nil, fmt.Errorf("Database Error Occured. [%s]", err.Error())
	}
	var result = s.mapper.LineToDto(*line)
	return &result, nil

}

func (s lineService) InsertLine(line *models.Line) (*models.Line, error) {
	return s.repo.Insert(line)
}

func (s lineService) GetLineList() ([]models.LineDto01, error) {
	return s.repo.LineList01()
}

func (s lineService) PostLine(line *models.Line) (*models.Line, error) {
	var selectedLine *models.Line = nil
	var err error = nil
	selectedLine, err = s.repo.SelectByCode(line.Line_Code)
	if err != nil {
		return nil, err
	}
	isNew := selectedLine == nil
	if isNew {
		return s.repo.Insert(line)
	} else {
		line.Id = selectedLine.Id
		return s.repo.Update(line)
	}
}

func (s lineService) WithTrx(trxHandle *gorm.DB) lineService {
	s.repo = s.repo.WithTx(trxHandle)
	return s
}

func (s lineService) PostLineArray(ctx context.Context, lines []models.Line) ([]models.Line, error) {
	var response []models.Line = make([]models.Line, 0)
	var trx = ctx.Value("db_tx").(*gorm.DB).Begin()
	for _, line := range lines {
		result, err := s.WithTrx(trx).PostLine(&line)
		if err != nil {
			trx.Rollback()
			return nil, err
		}
		response = append(response, *result)
	}
	if err := trx.Commit().Error; err != nil {
		return nil, err
	}
	return response, nil
}

func (s lineService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s lineService) InsertArray(entityArr []models.Line) ([]models.Line, error) {
	return s.repo.InsertArray(entityArr)
}

func (s lineService) InsertChunkArray(chunkSize int, allData []models.Line) error {
	// var maxSize = 1000
	var stratIndex = 0
	var endIndex = chunkSize
	if chunkSize > len(allData) {
		endIndex = len(allData) - 1
	}
	// txt := s.dbConnection.Begin()
	for {
		_, err := s.InsertArray(allData[stratIndex:endIndex])
		if err != nil {
			return err
		}
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
