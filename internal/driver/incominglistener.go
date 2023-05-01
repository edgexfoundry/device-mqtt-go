// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v3/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"
)

func (d *Driver) onIncomingDataReceived(_ mqtt.Client, message mqtt.Message) {
	var deviceName string
	var sourceName string

	var deviceResources []models.DeviceResource
	incomingTopic := message.Topic()
	subscribedTopic := d.serviceConfig.MQTTBrokerInfo.IncomingTopic
	subscribedTopic = strings.Replace(subscribedTopic, "#", "", -1)
	incomingTopic = strings.Replace(incomingTopic, subscribedTopic, "", -1)
	metaData := strings.Split(incomingTopic, "/")
	if len(metaData) != 2 {
		driver.Logger.Errorf("[Incoming listener] Incoming data ignored, incoming topic data should have format .../<device_name>/<source_name>: `%s`", incomingTopic)
		return
	}
	deviceName = metaData[0]
	sourceName = metaData[1]

	deviceCommand, ok := d.sdk.DeviceCommand(deviceName, sourceName)
	if !ok {
		deviceResource, ok := d.sdk.DeviceResource(deviceName, sourceName)
		if !ok {
			driver.Logger.Errorf("[Incoming listener] Incoming data ignored, source name `%s` not found as Device Command or Device Resource on the device `%s`", sourceName, deviceName)
			return
		}

		if !strings.Contains(strings.ToUpper(deviceResource.Properties.ReadWrite), "R") {
			driver.Logger.Errorf("[Incoming listener] Incoming data ignored, Device Resource `%s` not Readable", sourceName)
			return
		}

		deviceResources = append(deviceResources, deviceResource)
	} else {
		if !strings.Contains(strings.ToUpper(deviceCommand.ReadWrite), "R") {
			driver.Logger.Errorf("[Incoming listener] Incoming data ignored, Device Command `%s` not Readable", sourceName)
			return
		}

		for _, resourceOperation := range deviceCommand.ResourceOperations {
			deviceResource, ok := d.sdk.DeviceResource(deviceName, resourceOperation.DeviceResource)
			if !ok {
				driver.Logger.Errorf("[Incoming listener] Incoming data ignored, resource name `%s` from Device Command %s not found on the device `%s`", resourceOperation.DeviceResource, sourceName, deviceName)
				return
			}

			deviceResources = append(deviceResources, deviceResource)
		}
	}

	var asyncData map[string]interface{}
	err := json.Unmarshal(message.Payload(), &asyncData)
	if err != nil {
		driver.Logger.Errorf("[Incoming listener] Error un-marshaling incoming data : %v", err)
		return
	}

	var commandValues []*sdkModels.CommandValue

	for _, resource := range deviceResources {
		asyncValue, ok := asyncData[resource.Name]
		if !ok {
			driver.Logger.Errorf("[Incoming listener] Incoming data ignored: Resource Name %s not found in payload (%s)", resource.Name, string(message.Payload()))
			return
		}

		commandValue, err := newResult(resource, asyncValue)
		if err != nil {
			driver.Logger.Errorf("[Incoming listener] Incoming data ignored: %v", err)
			return
		}

		commandValues = append(commandValues, commandValue)

	}

	asyncValues := &sdkModels.AsyncValues{
		DeviceName:    deviceName,
		SourceName:    sourceName,
		CommandValues: commandValues,
	}

	driver.Logger.Debugf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload()))

	driver.AsyncCh <- asyncValues
}
