package mapper

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/cs161079/monorepo/common/utils"
	"github.com/fatih/structs"
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
				if err != nil {
					panic(err.Error())
				}
				v.Set(reflect.ValueOf(timeStr))
			case reflect.String.String():
				v.SetString(sourceFieldVal.(string))
			case reflect.Int64.String():
				val, _ := utils.StrToInt64(sourceFieldVal)
				v.Set(reflect.ValueOf(*val))
			case reflect.Int32.String():
				val, _ := utils.StrToInt32(sourceFieldVal)
				v.Set(reflect.ValueOf(*val))
			case reflect.Int16.String():
				val, _ := utils.StrToInt16(sourceFieldVal)
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

// Function to Map structures from one to another with same field data types
// but one of them has less fields from the other
func structMapper(source interface{}, target interface{}) {
	sourceMap := structs.Map(source)
	rvTarget := reflect.ValueOf(target)
	trvTarget := reflect.TypeOf(target)

	if rvTarget.Kind() == reflect.Pointer {
		rvTarget = rvTarget.Elem()
		trvTarget = trvTarget.Elem()
		target = reflect.New(rvTarget.Type())
	}
	for i := 0; i < rvTarget.NumField(); i++ {
		v := rvTarget.Field(i)
		tv := trvTarget.Field(i)
		// v.Set(reflect.ValueOf(source[tag]))
		fieldName := tv.Name
		sourceFieldVal := sourceMap[fieldName]
		if sourceFieldVal != nil {
			v.Set(reflect.ValueOf(sourceFieldVal))
		}
	}
}

func structMapper02(source any, target any) {
	sourceMap := structs.Map(source)
	// Convert map to JSON
	jsonData, err := json.Marshal(sourceMap)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		panic(err.Error())
	}
}
