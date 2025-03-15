package mapper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
)

const (
	TIME_TYPE = "time"
)

// With this Mapper map the records from the Oasa Server to structures that we have defined
// to implement procedures for the needs of the application
func internalMapper(source map[string]interface{}, target interface{}) {
	rvTarget := reflect.ValueOf(target)
	trvTarget := reflect.TypeOf(target)

	if rvTarget.Kind() == reflect.Pointer {
		rvTarget = rvTarget.Elem()
		trvTarget = trvTarget.Elem()
		target = reflect.New(rvTarget.Type())
	}
	for i := 0; i < rvTarget.NumField(); i++ {
		field := rvTarget.Field(i)
		fieldType := field.Kind().String()
		v := rvTarget.Field(i)
		tag, tagOk := trvTarget.Field(i).Tag.Lookup("oasa")
		if !tagOk || len(tag) == 0 {
			return
		}
		tagType, _ := trvTarget.Field(i).Tag.Lookup("type")
		if len(tagType) > 0 {
			fieldType = tagType
		}
		// v.Set(reflect.ValueOf(source[tag]))
		sourceFieldVal := source[tag]
		if sourceFieldVal != nil {
			switch fieldType {
			case TIME_TYPE:
				timeStamp, err := time.ParseInLocation("2006-01-02 15:04:05", sourceFieldVal.(string), time.Now().Location())
				if err != nil {
					panic(err.Error())
				}
				var timeStr = timeStamp.Format("15:04")
				v.Set(reflect.ValueOf(timeStr))
			case reflect.String.String():
				v.SetString(sourceFieldVal.(string))
			case reflect.Int64.String():
				val, _ := utils.StrToInt64(sourceFieldVal)
				v.Set(reflect.ValueOf(*val))
			case reflect.Int32.String():
				val, err := utils.StrToInt32(sourceFieldVal)
				if err != nil {
					logger.WARN(fmt.Sprintf("ERROR ON CONVERT STR TO INT32 %s", err.Error()))
				}
				v.Set(reflect.ValueOf(*val))
			case reflect.Int16.String():
				val, err := utils.StrToInt16(sourceFieldVal)
				if err != nil {
					logger.WARN(fmt.Sprintf("ERROR ON CONVERT STR TO INT16 %s", err.Error()))
				}
				v.Set(reflect.ValueOf(*val))
			case reflect.Int8.String():
				val, _ := utils.StrToInt8(sourceFieldVal)
				v.Set(reflect.ValueOf(*val))
			case reflect.Float32.String():
				v.Set(reflect.ValueOf(utils.StrToFloat32(sourceFieldVal)))
			case reflect.Float64.String():
				v.Set(reflect.ValueOf(utils.StrToFloat(sourceFieldVal)))
			case reflect.Ptr.String():
				v.Set(reflect.ValueOf(nil))
			}
		}

	}
}

// This function map one struct with other
// Convert source struct to JSON and make JSON as target struct.
//
//	@param source is any struct to read
//	@param target is any pointer of struct to write in
func structMapper02(source any, target any) {
	jsonData, err := json.Marshal(source)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		logger.ERROR(err.Error())
		panic(err.Error())
	}
}

// This function call structMapper02 in
// Created because structMapper02 is not accessible from other packages.
func MapStruct(source any, target any) {
	structMapper02(source, target)
}
