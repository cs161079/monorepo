package mapper

import (
	"encoding/json"
	"fmt"

	models "github.com/cs161079/monorepo/common/models"
)

type Schedule01Mapper interface {
	DtoToScheduleTime(source map[string]interface{}) (models.Scheduletime01Dto, error)
	ScheduletimeDtoToScheduletime(source models.ScheduletimeDto, direction int) models.Scheduletime
}

type schedule01Mapper struct {
}

func NewScheduletimeMapper() Schedule01Mapper {
	return schedule01Mapper{}
}

func (m schedule01Mapper) DtoToScheduleTime(source map[string]interface{}) (models.Scheduletime01Dto, error) {
	_, ok1 := source["go"]
	_, ok2 := source["come"]

	var result models.Scheduletime01Dto = models.Scheduletime01Dto{}

	if !ok1 || !ok2 {
		return result, fmt.Errorf("Το record είναι ελλειπής. %v", source)
	}

	bts, err := json.Marshal(source)
	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Σφάλμα κατά την μετατροπή από JSON. [%s]%v", err.Error(), source))
	}
	err = json.Unmarshal(bts, &result)
	if err != nil {
		return result, fmt.Errorf(fmt.Sprintf("Σφάλμα κατά την μετατροπή από JSON σε record. [%s] %v", err.Error(), source))
	}
	return result, nil
}

func (m schedule01Mapper) ScheduletimeDtoToScheduletime(source models.ScheduletimeDto, direction int) models.Scheduletime {
	var result models.Scheduletime = models.Scheduletime{}
	structMapper02(source, &result)
	result.Direction = int8(direction)
	if result.Direction == models.Direction_GO {
		result.Start_time = source.Start_time1
		result.End_time = source.End_time1
	} else {
		result.Start_time = source.Start_time2
		result.End_time = source.End_time2
	}

	return result
}
