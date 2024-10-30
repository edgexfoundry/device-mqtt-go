// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"math"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/spf13/cast"
)

// checkValueInRange checks value range is valid
func checkValueInRange(valueType string, reading interface{}) bool {
	isValid := false

	if valueType == common.ValueTypeString || valueType == common.ValueTypeBool || valueType == common.ValueTypeObject {
		return true
	}

	if valueType == common.ValueTypeInt8 || valueType == common.ValueTypeInt16 ||
		valueType == common.ValueTypeInt32 || valueType == common.ValueTypeInt64 {
		val := cast.ToInt64(reading)
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == common.ValueTypeUint8 || valueType == common.ValueTypeUint16 ||
		valueType == common.ValueTypeUint32 || valueType == common.ValueTypeUint64 {
		val := cast.ToUint64(reading)
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == common.ValueTypeFloat32 || valueType == common.ValueTypeFloat64 {
		val := cast.ToFloat64(reading)
		isValid = checkFloatValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType string, val uint64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeUint8:
		if val <= math.MaxUint8 {
			isValid = true
		}
	case common.ValueTypeUint16:
		if val <= math.MaxUint16 {
			isValid = true
		}
	case common.ValueTypeUint32:
		if val <= math.MaxUint32 {
			isValid = true
		}
	case common.ValueTypeUint64:
		maxiMum := uint64(math.MaxUint64)
		if val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType string, val int64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeInt8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case common.ValueTypeInt16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case common.ValueTypeInt32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case common.ValueTypeInt64:
		isValid = true
	}
	return isValid
}

func checkFloatValueRange(valueType string, val float64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeFloat32:
		if !math.IsNaN(val) && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case common.ValueTypeFloat64:
		if !math.IsNaN(val) && !math.IsInf(val, 0) {
			isValid = true
		}
	}
	return isValid
}
