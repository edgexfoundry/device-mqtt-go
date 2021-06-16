// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"

	"github.com/eclipse/paho.mqtt.golang"
)

func startCommandResponseListening() error {
	var qos = byte(driver.serviceConfig.MQTTBrokerInfo.Qos)
	var topic = driver.serviceConfig.MQTTBrokerInfo.ResponseTopic

	token := driver.mqttClient.Subscribe(topic, qos, onCommandResponseReceived)
	if token.Wait() && token.Error() != nil {
		driver.Logger.Infof("[Response listener] Stop command response listening. Cause:%v", token.Error())
		return token.Error()
	}

	driver.Logger.Info("[Response listener] Start command response listening. ")
	select {}
}

func onCommandResponseReceived(client mqtt.Client, message mqtt.Message) {
	var response map[string]interface{}

	json.Unmarshal(message.Payload(), &response)
	uuid, ok := response["uuid"].(string)
	if ok {
		driver.CommandResponses.Store(uuid, string(message.Payload()))
		driver.Logger.Debugf("[Response listener] Command response received: topic=%v uuid=%v msg=%v", message.Topic(), uuid, string(message.Payload()))
	} else {
		driver.Logger.Warnf("[Response listener] Command response ignored. No UUID found in the message: topic=%v msg=%v", message.Topic(), string(message.Payload()))
	}
}
