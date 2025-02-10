package mapper

import (
	models "github.com/cs161079/monorepo/common/models"
)

func NewLineMapper() LineMapper {
	return lineMapper{}
}

type LineMapper interface {
	GenDtLineOasa(map[string]interface{}) models.LineOasa
	OasaToLine(models.LineOasa) models.Line
	LineToDto(models.Line) models.LineDto
}

type lineMapper struct {
}

//	With this function we convert OASA data (JSON) into our own structures (Struct)
//
// @param Mapper key: string value: interface{}
// @return models.LineOasa
func (m lineMapper) GenDtLineOasa(source map[string]interface{}) models.LineOasa {
	var busLineOb models.LineOasa
	internalMapper(source, &busLineOb)

	return busLineOb
}

//	With this function we convert LineOasa Data into our own struct Line
//
// @param models.LineOasa
// @return models.Line
func (m lineMapper) OasaToLine(source models.LineOasa) models.Line {
	var target models.Line
	structMapper02(source, &target)
	return target
}

func (m lineMapper) LineToDto(source models.Line) models.LineDto {
	var result models.LineDto
	MapStruct(source, &result)
	return result
}
