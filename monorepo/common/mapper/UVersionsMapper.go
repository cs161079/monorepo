package mapper

import "github.com/cs161079/monorepo/common/models"

type UversionMapper interface {
	GeneralUVersions(any) models.UVersionsOasa
	OasaToUVersions(models.UVersionsOasa) models.UVersions
}

type uversionMapper struct{}

func NewUVersionMapper() UversionMapper {
	return uversionMapper{}
}
func (m uversionMapper) GeneralUVersions(source any) models.UVersionsOasa {
	var oasaOb models.UVersionsOasa
	vMap, ok := source.(map[string]interface{})
	if !ok {
		panic("An error occurred parsing the object.[Uversions Object from OASA]")
	}
	internalMapper(vMap, &oasaOb)

	return oasaOb
}

func (m uversionMapper) OasaToUVersions(source models.UVersionsOasa) models.UVersions {
	var target models.UVersions
	structMapper02(source, &target)
	return target
}
