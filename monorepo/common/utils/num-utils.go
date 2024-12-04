package utils

import (
	"strconv"
	"strings"
)

func StrToInt8(input interface{}) (*int8, error) {
	res, err := stringToNumberInternal(input.(string), 8)
	if err != nil {
		return nil, err
	}
	res1 := int8(*res)
	return &res1, err
}

func StrToInt64(input interface{}) (*int64, error) {
	res, err := stringToNumberInternal(input.(string), 64)
	if err != nil {
		return nil, err
	}
	res1 := int64(*res)
	return &res1, err
}

func StrToInt32(input interface{}) (*int32, error) {
	res, err := stringToNumberInternal(input.(string), 32)
	if err != nil {
		return nil, err
	}
	res1 := int32(*res)
	return &res1, err
}

func StrToInt16(input interface{}) (*int16, error) {
	res, err := stringToNumberInternal(input.(string), 16)
	if err != nil {
		return nil, err
	}
	res1 := int16(*res)
	return &res1, err
}

func stringToNumberInternal(input string, bitSize int) (*int64, error) {
	sourceNumVal, err := strconv.ParseInt(strings.Trim(input, " "), 10, bitSize)
	if err != nil {
		//panic("Δεν ήταν δυνατή η μετατροπή της συμβολοσειράς σε αριθμό για το πεδίο " + input)
		return nil, err
	}
	return &sourceNumVal, nil
}

func stringToFloatInternal(input string, bitSize int) float64 {
	sourceNumVal, error := strconv.ParseFloat(strings.Trim(input, " "), bitSize)
	if error != nil {
		panic("Δεν ήταν δυνατή η μετατροπή της συμβολοσειράς σε αριθμό για το πεδίο " + input)
	}
	return sourceNumVal
}

func StrToFloat(input interface{}) float64 {
	return stringToFloatInternal(input.(string), 64)
}
func StrToFloat32(input interface{}) float32 {
	return float32(stringToFloatInternal(input.(string), 32))
}
