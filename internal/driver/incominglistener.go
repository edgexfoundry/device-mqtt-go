// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	name = "name"
	cmd  = "cmd"
)

func (d *Driver) onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	var deviceName string
	var resourceName string
	var reading interface{}

	if d.serviceConfig.MQTTBrokerInfo.UseTopicLevels {
		incomingTopic := message.Topic()
		subscribedTopic := d.serviceConfig.MQTTBrokerInfo.IncomingTopic
		subscribedTopic = strings.Replace(subscribedTopic, "#", "", -1)
		incomingTopic = strings.Replace(incomingTopic, subscribedTopic, "", -1)
		metaData := strings.Split(incomingTopic, "/")
		if len(metaData) != 2 {
			driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, incoming topic data should have format .../<device_name>/<resource_name>: `%s`", incomingTopic)
			return
		}
		deviceName = metaData[0]
		resourceName = metaData[1]
		reading = string(message.Payload())
	} else {
		var data map[string]interface{}
		err := json.Unmarshal(message.Payload(), &data)
		if err != nil {
			driver.Logger.Errorf("Error unmarshaling payload: %s", err)
			return
		}
		nameVal, ok := data[name]
		if !ok {
			driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, reading data `%v` should contain the field `%s` to indicate the device name", data, name)
			return
		}
		cmdVal, ok := data[cmd]
		if !ok {
			driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, reading data `%v` should contain the field `%s` to indicate the device resource name", data, cmd)
			return
		}
		deviceName = fmt.Sprintf("%s", nameVal)
		resourceName = fmt.Sprintf("%s", cmdVal)

		reading, ok = data[resourceName]
		if !ok {
			driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, reading data `%v` should contain the field `%s` with the actual reading value", data, resourceName)
			return
		}
	}

	service := service.RunningService()

	deviceObject, ok := service.DeviceResource(deviceName, resourceName)
	if !ok {
		driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, device resource `%s` not found from the device `%s`", resourceName, deviceName)
		return
	}

	req := models.CommandRequest{
		DeviceResourceName: resourceName,
		Type:               deviceObject.Properties.ValueType,
	}

	result, err := newResult(req, reading)
	if err != nil {
		driver.Logger.Errorf("[Incoming listener] Incoming reading ignored, %v", err)
		return
	}

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{result},
	}

	driver.Logger.Debugf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload()))

	driver.AsyncCh <- asyncValues
}
