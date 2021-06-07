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
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

func startCommandResponseListening() error {

	var scheme = driver.serviceConfig.MQTTBrokerInfo.ResponseSchema
	var brokerUrl = driver.serviceConfig.MQTTBrokerInfo.ResponseHost
	var brokerPort = driver.serviceConfig.MQTTBrokerInfo.ResponsePort
	var authMode = driver.serviceConfig.MQTTBrokerInfo.ResponseAuthMode
	var secretPath = driver.serviceConfig.MQTTBrokerInfo.ResponseCredentialsPath
	var mqttClientId = driver.serviceConfig.MQTTBrokerInfo.ResponseClientId
	var qos = byte(driver.serviceConfig.MQTTBrokerInfo.ResponseQos)
	var keepAlive = driver.serviceConfig.MQTTBrokerInfo.ResponseKeepAlive
	var topic = driver.serviceConfig.MQTTBrokerInfo.ResponseTopic

	uri := &url.URL{
		Scheme: strings.ToLower(scheme),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
	}

	err := SetCredentials(uri, "Response", authMode, secretPath)
	if err != nil {
		return err
	}

	var client mqtt.Client
	for i := 1; i <= driver.serviceConfig.MQTTBrokerInfo.ConnEstablishingRetry; i++ {
		client, err = createClient(mqttClientId, uri, keepAlive)
		if err != nil && i == driver.serviceConfig.MQTTBrokerInfo.ConnEstablishingRetry {
			return err
		} else if err != nil {
			driver.Logger.Error(fmt.Sprintf("Fail to initial conn for command response, %v ", err))
			time.Sleep(time.Duration(driver.serviceConfig.MQTTBrokerInfo.ConnEstablishingRetry) * time.Second)
			driver.Logger.Warn("Retry to initial conn for command response")
			continue
		}
		break
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
