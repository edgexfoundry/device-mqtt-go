// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go/v3/pkg/models"
)

const (
	name = "name"
	cmd  = "cmd"
)

func (d *Driver) onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	var deviceName string
	var resourceName string
	var reading interface{}

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

	deviceObject, ok := d.sdk.DeviceResource(deviceName, resourceName)
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
