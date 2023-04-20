// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (d *Driver) onCommandResponseReceived(client mqtt.Client, message mqtt.Message) {
	var uuid string

	topic := message.Topic()
	metaData := strings.Split(topic, "/")

	if len(metaData) == 0 {
		driver.Logger.Errorf("[Response listener] Command response ignored. metaData in the message is not sufficient to retrieve UUID: topic=%v msg=%v", message.Topic(), metaData)
		return
	} else {
		uuid = metaData[len(metaData)-1]
	}

	driver.CommandResponses.Store(uuid, string(message.Payload()))
	driver.Logger.Debugf("[Response listener] Command response received: topic=%v uuid=%v msg=%v", message.Topic(), uuid, string(message.Payload()))
}
