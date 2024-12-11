package service

import (
	"context"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/repository"

	"gorm.io/gorm"
)

type LineService interface {
	SelectByLineCode(lineCode int32) (*models.Line, error)
	InsertLine(line *models.Line) (*models.Line, error)
	InsertArray([]models.Line) ([]models.Line, error)
	PostLine(line *models.Line) (*models.Line, error)
	PostLineArray(context.Context, []models.Line) ([]models.Line, error)
	WithTrx(*gorm.DB) lineService
	DeleteAll() error
	GetMapper() mapper.LineMapper
	GetLineList() ([]models.Line, error)
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

func (s lineService) SelectByLineCode(line_code int32) (*models.Line, error) {
	return s.repo.SelectByCode(line_code)

}

func (s lineService) InsertLine(line *models.Line) (*models.Line, error) {
	return s.repo.Insert(line)
}

func (s lineService) GetLineList() ([]models.Line, error) {
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
