// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"math"

	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/spf13/cast"
)

// checkValueInRange checks value range is valid
func checkValueInRange(valueType sdkModel.ValueType, reading interface{}) bool {
	isValid := false

	if valueType == sdkModel.String || valueType == sdkModel.Bool {
		return true
	}

	if valueType == sdkModel.Int8 || valueType == sdkModel.Int16 ||
		valueType == sdkModel.Int32 || valueType == sdkModel.Int64 {
		val := cast.ToInt64(reading)
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == sdkModel.Uint8 || valueType == sdkModel.Uint16 ||
		valueType == sdkModel.Uint32 || valueType == sdkModel.Uint64 {
		val := cast.ToUint64(reading)
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == sdkModel.Float32 || valueType == sdkModel.Float64 {
		val := cast.ToFloat64(reading)
		isValid = checkFloatValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType sdkModel.ValueType, val uint64) bool {
	var isValid = false
	switch valueType {
	case sdkModel.Uint8:
		if val >= 0 && val <= math.MaxUint8 {
			isValid = true
		}
	case sdkModel.Uint16:
		if val >= 0 && val <= math.MaxUint16 {
			isValid = true
		}
	case sdkModel.Uint32:
		if val >= 0 && val <= math.MaxUint32 {
			isValid = true
		}
	case sdkModel.Uint64:
		maxiMum := uint64(math.MaxUint64)
		if val >= 0 && val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType sdkModel.ValueType, val int64) bool {
	var isValid = false
	switch valueType {
	case sdkModel.Int8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case sdkModel.Int16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case sdkModel.Int32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case sdkModel.Int64:
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			isValid = true
		}
	}
	return isValid
}

func checkFloatValueRange(valueType sdkModel.ValueType, val float64) bool {
	var isValid = false
	switch valueType {
	case sdkModel.Float32:
		if math.Abs(val) >= math.SmallestNonzeroFloat32 && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case sdkModel.Float64:
		if math.Abs(val) >= math.SmallestNonzeroFloat64 && math.Abs(val) <= math.MaxFloat64 {
			isValid = true
		}
	}
	return isValid
}
