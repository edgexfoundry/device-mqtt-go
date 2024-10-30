// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	driver = new(Driver)
	driver.Logger = logger.NewMockClient()
}

func TestNewResult_bool(t *testing.T) {
	var reading interface{} = true
	resource := models.DeviceResource{
		Name:       "light",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeBool},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.BoolValue()
	if val != true || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint8(t *testing.T) {
	var reading interface{} = uint8(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int8(t *testing.T) {
	var reading interface{} = int8(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResultFailed_int8(t *testing.T) {
	var reading interface{} = int16(256)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt8},
	}

	_, err := newResult(resource, reading)
	fmt.Println(err)
	if err == nil || !strings.Contains(err.Error(), "Reading 256 is out of the value type(Int8)'s range") {
		t.Errorf("Convert new result should be failed")
	}
}

func TestNewResult_uint16(t *testing.T) {
	var reading interface{} = uint16(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint16},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint16Value()
	if val != uint16(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int16(t *testing.T) {
	var reading interface{} = int16(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt16},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int16Value()
	if val != int16(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint32(t *testing.T) {
	var reading interface{} = uint32(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int32(t *testing.T) {
	var reading interface{} = int32(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int32Value()
	if val != int32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint64(t *testing.T) {
	var reading interface{} = uint64(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int64(t *testing.T) {
	var reading interface{} = int64(123)
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float32(t *testing.T) {
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeFloat32},
	}

	tests := []struct {
		name     string
		resource models.DeviceResource
		reading  interface{}
		expected interface{}
	}{
		{"Valid string 0", resource, "0", float32(0)},
		{"Valid string 123.32", resource, "123.32", float32(123.32)},
		{"Valid float32 0", resource, 0, float32(0)},
		{"Valid float32 123.32", resource, 123.32, float32(123.32)},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cmdVal, err := newResult(resource, testCase.reading)
			require.NoError(t, err)
			result, err := cmdVal.Float32Value()
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNewResult_float64(t *testing.T) {
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeFloat64},
	}

	tests := []struct {
		name     string
		resource models.DeviceResource
		reading  interface{}
		expected interface{}
	}{
		{"Valid string 0", resource, "0", float64(0)},
		{"Valid string 0.123", resource, "0.123", 0.123},
		{"Valid float64 0", resource, 0, float64(0)},
		{"Valid float64 0.123", resource, 0.123, 0.123},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cmdVal, err := newResult(resource, testCase.reading)
			require.NoError(t, err)
			result, err := cmdVal.Float64Value()
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNewResult_float64ToInt8(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt16(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt16},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int16Value()
	if val != int16(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt32(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int32Value()
	if val != int32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt64(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint8(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint16(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint16},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint16Value()
	if val != uint16(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint32(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint64(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToFloat32(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeFloat32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Float32Value()
	if val != float32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToString(t *testing.T) {
	var reading interface{} = 123.11
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeString},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != fmt.Sprintf("%v", reading) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_string(t *testing.T) {
	var reading interface{} = "test string"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeString},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != "test string" || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToInt64(t *testing.T) {
	var reading interface{} = "123"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToInt8(t *testing.T) {
	var reading interface{} = "123"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeInt8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint8(t *testing.T) {
	var reading interface{} = "123"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint8},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint32(t *testing.T) {
	var reading interface{} = "123"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint64(t *testing.T) {
	var reading interface{} = "123"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToBool(t *testing.T) {
	var reading interface{} = "true"
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeBool},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.BoolValue()
	if val != true || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_numberToUint64(t *testing.T) {
	var reading interface{} = 123
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeUint64},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_floatNumberToFloat32(t *testing.T) {
	var reading interface{} = 123.0
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeFloat32},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Float32Value()
	if val != float32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_numberToString(t *testing.T) {
	var reading interface{} = 123
	resource := models.DeviceResource{
		Name:       "temperature",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeString},
	}

	cmdVal, err := newResult(resource, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != "123" || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_object(t *testing.T) {
	var reading interface{} = map[string]string{"json": "My JSON"}
	resource := models.DeviceResource{
		Name:       "json",
		Properties: models.ResourceProperties{ValueType: common.ValueTypeObject},
	}

	cmdVal, err := newResult(resource, reading)
	require.NoError(t, err)
	val, err := cmdVal.ObjectValue()
	require.NoError(t, err)
	assert.Equal(t, reading, val)
}
