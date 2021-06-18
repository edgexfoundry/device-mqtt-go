// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"

	"github.com/eclipse/paho.mqtt.golang"
)

func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	var data map[string]interface{}
	json.Unmarshal(message.Payload(), &data)

	if !checkDataWithKey(data, "name") || !checkDataWithKey(data, "cmd") {
		return
	}

	deviceName := data["name"].(string)
	cmd := data["cmd"].(string)

	reading, ok := data[cmd]
	if !ok {
		driver.Logger.Warnf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload()))
		return
	}

	service := service.RunningService()

	deviceObject, ok := service.DeviceResource(deviceName, cmd)
	if !ok {
		driver.Logger.Warnf("[Incoming listener] Incoming reading ignored. No DeviceObject found : topic=%v msg=%v", message.Topic(), string(message.Payload()))
		return
	}

	req := models.CommandRequest{
		DeviceResourceName: cmd,
		Type:               deviceObject.Properties.ValueType,
	}

	result, err := newResult(req, reading)

	if err != nil {
		driver.Logger.Warnf("[Incoming listener] Incoming reading ignored.   topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err)
		return
	}

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{result},
	}

	driver.Logger.Debugf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload()))

	driver.AsyncCh <- asyncValues

}

func checkDataWithKey(data map[string]interface{}, key string) bool {
	val, ok := data[key]
	if !ok {
		driver.Logger.Warnf("[Incoming listener] Incoming reading ignored. No %v found : msg=%v", key, data)
		return false
	}

	switch val.(type) {
	case string:
		return true
	default:
		driver.Logger.Warnf("[Incoming listener] Incoming reading ignored. %v should be string : msg=%v", key, data)
		return false
	}
}
