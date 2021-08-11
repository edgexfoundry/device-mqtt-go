// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
)

func (d *Driver) onCommandResponseReceived(client mqtt.Client, message mqtt.Message) {
	var uuid string

	if d.serviceConfig.MQTTBrokerInfo.UseTopicLevels {
		topic := message.Topic()
		metaData := strings.Split(topic, "/")
		uuid = metaData[len(metaData)-1]
	} else {
		var response map[string]interface{}
		var ok bool

		json.Unmarshal(message.Payload(), &response)
		uuid, ok = response["uuid"].(string)
		if !ok {
			driver.Logger.Warnf("[Response listener] Command response ignored. No UUID found in the message: topic=%v msg=%v", message.Topic(), string(message.Payload()))
			return
		}
	}
	driver.CommandResponses.Store(uuid, string(message.Payload()))
	driver.Logger.Debugf("[Response listener] Command response received: topic=%v uuid=%v msg=%v", message.Topic(), uuid, string(message.Payload()))
}
