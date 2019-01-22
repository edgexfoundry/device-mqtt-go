package driver

import (
	"math"
	"strings"
)

// checkValueInRange checks value range is valid
func checkValueInRange(profileValueType string, readingValueType string, reading interface{}) bool {
	isValid := false

	if profileValueType == "string" || profileValueType == "bool" {
		return true
	}

	switch {
	case strings.Contains(profileValueType, "uint"):
		var val uint64
		if readingValueType == "int" {
			val = uint64(reading.(int64))
		} else if readingValueType == "uint" {
			val = uint64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = uint64(reading.(float32))
		} else {
			val = uint64(reading.(float64))
		}
		isValid = checkUintValueRange(profileValueType, val)

	case strings.Contains(profileValueType, "int"):
		var val int64
		if readingValueType == "int" {
			val = int64(reading.(int64))
		} else if readingValueType == "uint" {
			val = int64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = int64(reading.(float32))
		} else {
			val = int64(reading.(float64))
		}
		isValid = checkIntValueRange(profileValueType, val)

	case strings.Contains(profileValueType, "float"):
		var val float64
		if readingValueType == "int" {
			val = float64(reading.(int64))
		} else if readingValueType == "uint" {
			val = float64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = float64(reading.(float32))
		} else {
			val = float64(reading.(float64))
		}
		isValid = checkFloatValueRange(profileValueType, val)
	}

	return isValid
}

func checkUintValueRange(profileValueType string, val uint64) bool {
	var isValid = false
	switch profileValueType {
	case "uint8":
		if val >= 0 && val <= math.MaxUint8 {
			isValid = true
		}
	case "uint16":
		if val >= 0 && val <= math.MaxUint16 {
			isValid = true
		}
	case "uint32":
		if val >= 0 && val <= math.MaxUint32 {
			isValid = true
		}
	case "uint64":
		maxiMum := uint64(math.MaxUint64)
		if val >= 0 && val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(profileValueType string, val int64) bool {
	var isValid = false
	switch profileValueType {
	case "int8":
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case "int16":
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case "int32":
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case "int64":
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			isValid = true
		}
	}
	return isValid
}

func checkFloatValueRange(profileValueType string, val float64) bool {
	var isValid = false
	switch profileValueType {
	case "float32":
		if val >= math.SmallestNonzeroFloat32 && val <= math.MaxFloat32 {
			isValid = true
		}
	case "float64":
		if val >= math.SmallestNonzeroFloat64 && val <= math.MaxFloat64 {
			isValid = true
		}
	}
	return isValid
}
