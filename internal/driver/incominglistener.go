// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"github.com/cisco/senml"
	"net/url"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	sdk "github.com/edgexfoundry/device-sdk-go"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
)

type MqttReading struct {
	// Device Name received
	DeviceName string
	// Operation received
	Operation string
	// Raw value received
	RawValue interface{}
	// Unit of the value received
	Unit string
}

func startIncomingListening() error {
	var scheme = driver.Config.Incoming.Protocol
	var brokerUrl = driver.Config.Incoming.Host
	var brokerPort = driver.Config.Incoming.Port
	var username = driver.Config.Incoming.Username
	var password = driver.Config.Incoming.Password
	var mqttClientId = driver.Config.Incoming.MqttClientId
	var qos = byte(driver.Config.Incoming.Qos)
	var keepAlive = driver.Config.Incoming.KeepAlive
	var topic = driver.Config.Incoming.Topic

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

	token := client.Subscribe(topic, qos, onIncomingDataReceived)
	if token.Wait() && token.Error() != nil {
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		return token.Error()
	}

	driver.Logger.Info("[Incoming listener] Start incoming data listening. ")
	select {}
}


func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {

	rawSenML, err := senml.Decode(message.Payload(), senml.JSON)

	// check if we have received data in SenML format
	if err == nil {
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming data: SenML data received: topic=%v msg=%v", message.Topic(), string(message.Payload())))
		onIncomingSenMLDataReceived(rawSenML)
	} else {
		driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming data: Regular JSON data received: topic=%v msg=%v", message.Topic(), string(message.Payload())))
		onRegularIncomingDataReceived(message)
	}
}

func onIncomingSenMLDataReceived(message senml.SenML) {

	reading := MqttReading{}

	for _,v := range message.Records {
		reading.DeviceName= v.BaseName
		reading.Operation = v.Name
		reading.Unit = v.Unit


		switch {
		case v.Value != nil:
			reading.RawValue = *v.Value
		case v.BoolValue != nil:
			reading.RawValue = v.BoolValue
		case v.DataValue != "":
			reading.RawValue = v.DataValue
		case v.StringValue != "":
			reading.RawValue = v.StringValue
		}


		executeIncomingDataReceived(reading)

	}
}

func onRegularIncomingDataReceived (message mqtt.Message) {

	var data map[string]interface{}
	json.Unmarshal(message.Payload(), &data)

	if !checkDataWithKey(data, "name") || !checkDataWithKey(data, "cmd") {
		return
	}

	reading := MqttReading{}

	reading.DeviceName = data["name"].(string)
	reading.Operation = data["cmd"].(string)

	readingValue, ok := data[reading.Operation]
	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No reading data found : topic=%v msg=%v", message.Topic(), string(message.Payload())))
		return
	}

	reading.RawValue = readingValue

	executeIncomingDataReceived(reading)
}

func executeIncomingDataReceived (reading MqttReading) {

	service := sdk.RunningService()

	deviceObject, ok := service.DeviceObject(reading.DeviceName, reading.Operation, "get")

	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No DeviceObject found"))
		return
	}

	ro, ok := service.ResourceOperation(reading.DeviceName, reading.Operation, "get")

	if !ok {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No ResourceOperation found"))
		return
	}

	result, err := newResult(deviceObject, ro, reading.RawValue)


	if err != nil {
		driver.Logger.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. error=%v", err))
		return
	}

	if reading.Unit != "" {
		result.Unit = reading.Unit
	} else {
		result.Unit = deviceObject.Properties.Units.DefaultValue
	}

	asyncValues := &sdkModel.AsyncValues{
		DeviceName:    reading.DeviceName,
		CommandValues: []*sdkModel.CommandValue{result},
	}

	driver.Logger.Info(fmt.Sprintf("[Incoming listener] Incoming reading received"))

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
