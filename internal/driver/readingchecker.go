// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"math"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2"
	"github.com/spf13/cast"
)

// checkValueInRange checks value range is valid
func checkValueInRange(valueType string, reading interface{}) bool {
	isValid := false

	if valueType == v2.ValueTypeString || valueType == v2.ValueTypeBool {
		return true
	}

	if valueType == v2.ValueTypeInt8 || valueType == v2.ValueTypeInt16 ||
		valueType == v2.ValueTypeInt32 || valueType == v2.ValueTypeInt64 {
		val := cast.ToInt64(reading)
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == v2.ValueTypeUint8 || valueType == v2.ValueTypeUint16 ||
		valueType == v2.ValueTypeUint32 || valueType == v2.ValueTypeUint64 {
		val := cast.ToUint64(reading)
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == v2.ValueTypeFloat32 || valueType == v2.ValueTypeFloat64 {
		val := cast.ToFloat64(reading)
		isValid = checkFloatValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType string, val uint64) bool {
	var isValid = false
	switch valueType {
	case v2.ValueTypeUint8:
		if val >= 0 && val <= math.MaxUint8 {
			isValid = true
		}
	case v2.ValueTypeUint16:
		if val >= 0 && val <= math.MaxUint16 {
			isValid = true
		}
	case v2.ValueTypeUint32:
		if val >= 0 && val <= math.MaxUint32 {
			isValid = true
		}
	case v2.ValueTypeUint64:
		maxiMum := uint64(math.MaxUint64)
		if val >= 0 && val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType string, val int64) bool {
	var isValid = false
	switch valueType {
	case v2.ValueTypeInt8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case v2.ValueTypeInt16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case v2.ValueTypeInt32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case v2.ValueTypeInt64:
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			isValid = true
		}
	}
	return isValid
}

func checkFloatValueRange(valueType string, val float64) bool {
	var isValid = false
	switch valueType {
	case v2.ValueTypeFloat32:
		if !math.IsNaN(val) && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case v2.ValueTypeFloat64:
		if !math.IsNaN(val) && !math.IsInf(val, 0) {
			isValid = true
		}
	}
	return isValid
}
