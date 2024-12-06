package mapper

import (
	models "github.com/cs161079/monorepo/common/models"
)

func NewLineMapper() LineMapper {
	return lineMapper{}
}

type LineMapper interface {
	GeneralLine(map[string]interface{}) models.LineOasa
	OasaToLineDto(models.LineOasa) models.LineDto
	OasaToLine(models.LineOasa) models.Line
	LineDtoToLine(models.LineDto) models.Line
}

type lineMapper struct {
}

func (m lineMapper) GeneralLine(source map[string]interface{}) models.LineOasa {
	var busLineOb models.LineOasa
	internalMapper(source, &busLineOb)

	return busLineOb
}

func (m lineMapper) OasaToLineDto(source models.LineOasa) models.LineDto {
	var target models.LineDto
	structMapper02(source, &target)
	return target
}

func (m lineMapper) OasaToLine(source models.LineOasa) models.Line {
	var target models.Line
	structMapper02(source, &target)
	return target
}

func (m lineMapper) LineDtoToLine(source models.LineDto) models.Line {
	var target models.Line
	structMapper02(source, &target)
	return target
}
