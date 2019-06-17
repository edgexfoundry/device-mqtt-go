// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
)

func startCommandResponseListening() error {
	var scheme = driver.Config.ResponseSchema
	var brokerUrl = driver.Config.ResponseHost
	var brokerPort = driver.Config.ResponsePort
	var username = driver.Config.ResponseUser
	var password = driver.Config.ResponsePassword
	var mqttClientId = driver.Config.ResponseClientId
	var qos = byte(driver.Config.ResponseQos)
	var keepAlive = driver.Config.ResponseKeepAlive
	var topic = driver.Config.ResponseTopic

	uri := &url.URL{
		Scheme: strings.ToLower(scheme),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createClient(mqttClientId, uri, keepAlive)
	if err != nil {
		return err
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	token := client.Subscribe(topic, qos, onCommandResponseReceived)
	if token.Wait() && token.Error() != nil {
		driver.Logger.Info(fmt.Sprintf("[Response listener] Stop command response listening. Cause:%v", token.Error()))
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
		driver.Logger.Info(fmt.Sprintf("[Response listener] Command response received: topic=%v uuid=%v msg=%v", message.Topic(), uuid, string(message.Payload())))
	} else {
		driver.Logger.Warn(fmt.Sprintf("[Response listener] Command response ignored. No UUID found in the message: topic=%v msg=%v", message.Topic(), string(message.Payload())))
	}

}
