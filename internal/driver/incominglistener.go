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

	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func startIncomingListening() error {
	var scheme = driver.Config.IncomingSchema
	var brokerUrl = driver.Config.IncomingHost
	var brokerPort = driver.Config.IncomingPort
	var secretPath = driver.Config.IncomingCredentialsPath
	var mqttClientId = driver.Config.IncomingClientId
	var qos = byte(driver.Config.IncomingQos)
	var keepAlive = driver.Config.IncomingKeepAlive
	var topic = driver.Config.IncomingTopic

	credentials, err := GetCredentials(secretPath)
	if err != nil {
		return fmt.Errorf("Unable to get incoming MQTT credentials for secret path '%s': %s", secretPath, err.Error())
	}

	driver.Logger.Info("Incoming MQTT credentials loaded")

	uri := &url.URL{
		Scheme: strings.ToLower(scheme),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(credentials.Username, credentials.Password),
	}

	var client mqtt.Client
	for i := 1; i <= driver.Config.ConnEstablishingRetry; i++ {
		client, err = createClient(mqttClientId, uri, keepAlive)
		if err != nil && i == driver.Config.ConnEstablishingRetry {
			return err
		} else if err != nil {
			driver.Logger.Error(fmt.Sprintf("Fail to initial conn for incoming data, %v ", err))
			time.Sleep(time.Duration(driver.Config.ConnEstablishingRetry) * time.Second)
			driver.Logger.Warn("Retry to initial conn for incoming data")
			continue
		}
		break
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	token := client.Subscribe(topic, qos, onIncomingDataReceived)
	if token.Wait() && token.Error() != nil {
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		return token.Error()
	}

	driver.Logger.Info("[Incoming listener] Start incoming data listening. ")
	select {}
}

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
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	service := sdk.RunningService()

	deviceObject, ok := service.DeviceResource(deviceName, cmd, "get")
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	req := sdkModel.CommandRequest{
		DeviceResourceName: cmd,
		Type:               sdkModel.ParseValueType(deviceObject.Properties.Value.Type),
	}

	result, err := newResult(req, reading)

	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored.   topic=%v msg=%v error=%v", message.Topic(), string(message.Payload()), err))
		return
	}

	asyncValues := &sdkModel.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*sdkModel.CommandValue{result},
	}

	driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming reading received: topic=%v msg=%v", message.Topic(), string(message.Payload())))

	driver.AsyncCh <- asyncValues

}

func checkDataWithKey(data map[string]interface{}, key string) bool {
	val, ok := data[key]
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No %v found : msg=%v", key, data))
		return false
	}

	switch val.(type) {
	case string:
		return true
	default:
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. %v should be string : msg=%v", key, data))
		return false
	}
}
