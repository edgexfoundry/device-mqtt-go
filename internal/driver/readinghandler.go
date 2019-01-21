package driver

import (
	"fmt"
	"strconv"
	"strings"
)

// handleReadingStringValue checks and converts reading when it is string type
func handleReadingStringValue(profileValueType string, readingValueType string, reading interface{}) (interface{}, string, error) {
	if profileValueType == "string" {
		reading = fmt.Sprintf("%v", reading)
		readingValueType = "string"
		return reading, readingValueType, nil

	} else if readingValueType != "string" {
		return reading, readingValueType, nil
	}

	// Parse reading string according to profileValueType. Number will be convert to int or unit and then check number's value range in the next function.
	switch {
	case profileValueType == "bool":
		reading, err := strconv.ParseBool(reading.(string))
		readingValueType = "bool"
		if err != nil {
			err = fmt.Errorf("parse result fail. Reading's value (%v) can't parse to bool ,err:%v", reading, err)
			return reading, readingValueType, err
		} else {
			return reading, readingValueType, nil
		}
	case strings.Contains(profileValueType, "uint"):
		reading, err := strconv.ParseUint(reading.(string), 10, 64)
		readingValueType = "uint"
		if err != nil {
			err = fmt.Errorf("parse result fail. Reading's value (%v) can't parse to uint ,err:%v", reading, err)
			return reading, readingValueType, err
		} else {
			return reading, readingValueType, nil
		}
	case strings.Contains(profileValueType, "int"):
		reading, err := strconv.ParseInt(reading.(string), 10, 64)
		readingValueType = "int"
		if err != nil {
			err = fmt.Errorf("parse result fail. Reading's value (%v) can't parse to int ,err:%v", reading, err)
			return reading, readingValueType, err
		} else {
			return reading, readingValueType, nil
		}
	case profileValueType == "float32" || profileValueType == "float":
		val, err := strconv.ParseFloat(reading.(string), 32)
		readingValueType = "float32"
		if err != nil {
			err = fmt.Errorf("parse result fail. Reading's value (%v) can't parse to float32 ,err:%v", reading, err)
			return reading, readingValueType, err
		} else {
			reading = float32(val)
			return reading, readingValueType, nil
		}
	case profileValueType == "float64":
		reading, err := strconv.ParseFloat(reading.(string), 64)
		readingValueType = "float64"
		if err != nil {
			err = fmt.Errorf("parse result fail. Reading's value (%v) can't parse to float64 ,err:%v", reading, err)
			return reading, readingValueType, err
		} else {
			return reading, readingValueType, nil
		}
	default:
		return reading, readingValueType, nil
	}
}
