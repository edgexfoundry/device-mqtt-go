// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	logger "github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"gopkg.in/mgo.v2/bson"
)

var once sync.Once
var driver *Driver

type Config struct {
	Incoming connectionInfo
	Response connectionInfo
}

type connectionInfo struct {
	MqttProtocol   string
	MqttBroker     string
	MqttBrokerPort int
	MqttClientID   string
	MqttTopic      string
	MqttQos        int
	MqttUser       string
	MqttPassword   string
	MqttKeepAlive  int
}

type Driver struct {
	Logger           logger.LoggingClient
	AsyncCh          chan<- *sdkModel.AsyncValues
	CommandResponses map[string]string
	Config           *configuration
}

func NewProtocolDriver() sdkModel.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
		driver.CommandResponses = make(map[string]string)
	})
	return driver
}

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModel.AsyncValues) error {
	d.Logger = lc
	d.AsyncCh = asyncCh

	config, err := CreateDriverConfig(device.DriverConfigs())
	if err != nil {
		panic(fmt.Errorf("read MQTT driver configuration failed: %v", err))
	}
	d.Config = config

	go func() {
		err := startCommandResponseListening()
		if err != nil {
			panic(fmt.Errorf("start command response Listener failed, please check MQTT broker settings are correct, %v", err))
		}
	}()

	go func() {
		err := startIncomingListening()
		if err != nil {
			panic(fmt.Errorf("start incoming data Listener failed, please check MQTT broker settings are correct, %v", err))
		}
	}()

	return nil
}

func (d *Driver) DisconnectDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Warn("Driver's DisconnectDevice function didn't implement")
	return nil
}

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var err error

	// create device client and open connection
	connectionInfo, err := CreateConnectionInfo(protocols)
	if err != nil {
		return responses, err
	}

	uri := &url.URL{
		Scheme: strings.ToLower(connectionInfo.Schema),
		Host:   fmt.Sprintf("%s:%s", connectionInfo.Host, connectionInfo.Port),
		User:   url.UserPassword(connectionInfo.User, connectionInfo.Password),
	}

	client, err := createClient(connectionInfo.ClientId, uri, 30)
	if err != nil {
		return responses, err
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	for i, req := range reqs {
		res, err := d.handleReadCommandRequest(client, req, connectionInfo.Topic)
		if err != nil {
			driver.Logger.Info(fmt.Sprintf("Handle read commands failed: %v", err))
			return responses, err
		}

		responses[i] = res
	}

	return responses, err
}

func (d *Driver) handleReadCommandRequest(deviceClient MQTT.Client, req sdkModel.CommandRequest, topic string) (*sdkModel.CommandValue, error) {
	var result = &sdkModel.CommandValue{}
	var err error
	var qos = byte(0)
	var retained = false

	var method = "get"
	var cmdUuid = bson.NewObjectId().Hex()
	var cmd = req.DeviceResource.Name

	data := make(map[string]interface{})
	data["uuid"] = cmdUuid
	data["method"] = method
	data["cmd"] = cmd

	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	deviceClient.Publish(topic, qos, retained, jsonData)

	driver.Logger.Info(fmt.Sprintf("Publish command: %v", string(jsonData)))

	// fetch response from MQTT broker after publish command successful
	cmdResponse, ok := fetchCommandResponse(d.CommandResponses, cmdUuid)
	if !ok {
		err = fmt.Errorf("can not fetch command response: method=%v cmd=%v", method, cmd)
		return result, err
	}

	driver.Logger.Info(fmt.Sprintf("Parse command response: %v", cmdResponse))

	var response map[string]interface{}
	json.Unmarshal([]byte(cmdResponse), &response)
	reading, ok := response[req.DeviceResource.Name]
	if !ok {
		err = fmt.Errorf("can not fetch command reading: method=%v cmd=%v", method, cmd)
		return result, err
	}

	result, err = newResult(req.DeviceResource, req.RO, reading)
	if err != nil {
		return result, err
	} else {
		driver.Logger.Info(fmt.Sprintf("Get command finished: %v", result))
	}

	return result, err
}

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {
	var err error

	// create device client and open connection
	connectionInfo, err := CreateConnectionInfo(protocols)
	if err != nil {
		return err
	}

	uri := &url.URL{
		Scheme: strings.ToLower(connectionInfo.Schema),
		Host:   fmt.Sprintf("%s:%s", connectionInfo.Host, connectionInfo.Port),
		User:   url.UserPassword(connectionInfo.User, connectionInfo.Password),
	}

	client, err := createClient(connectionInfo.ClientId, uri, 30)
	if err != nil {
		return err
	}
	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	for i, req := range reqs {
		err = d.handleWriteCommandRequest(client, req, connectionInfo.Topic, params[i])
		if err != nil {
			driver.Logger.Info(fmt.Sprintf("Handle write commands failed: %v", err))
			return err
		}
	}

	return err
}

func (d *Driver) handleWriteCommandRequest(deviceClient MQTT.Client, req sdkModel.CommandRequest, topic string, param *sdkModel.CommandValue) error {
	var err error
	var qos = byte(0)
	var retained = false

	var method = "set"
	var cmdUuid = bson.NewObjectId().Hex()
	var cmd = req.DeviceResource.Name

	data := make(map[string]interface{})
	data["uuid"] = cmdUuid
	data["method"] = method
	data["cmd"] = cmd

	commandValue, err := newCommandValue(req.DeviceResource, param)
	if err != nil {
		return err
	} else {
		data[cmd] = commandValue
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	deviceClient.Publish(topic, qos, retained, jsonData)

	driver.Logger.Info(fmt.Sprintf("Publish command: %v", string(jsonData)))

	//wait and fetch response from CommandResponses map
	var cmdResponse string
	var ok bool
	for i := 0; i < 5; i++ {
		cmdResponse, ok = d.CommandResponses[cmdUuid]
		if ok {
			break
		} else {
			time.Sleep(time.Second * time.Duration(1))
		}
	}

	if !ok {
		err = fmt.Errorf("can not fetch command response: method=%v cmd=%v", method, cmd)
		return err
	}

	driver.Logger.Info(fmt.Sprintf("Put command finished: %v", cmdResponse))

	return nil
}

func (d *Driver) Stop(force bool) error {
	d.Logger.Warn("Driver's Stop function didn't implement")
	return nil
}

// Create a MQTT client
func createClient(clientID string, uri *url.URL, keepAlive int) (MQTT.Client, error) {
	driver.Logger.Info(fmt.Sprintf("Create MQTT client and connection: uri=%v clientID=%v ", uri.String(), clientID))
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetClientID(clientID)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetKeepAlive(time.Second * time.Duration(keepAlive))
	opts.SetConnectionLostHandler(func(client MQTT.Client, e error) {
		driver.Logger.Warn(fmt.Sprintf("Connection lost : %v", e))
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			driver.Logger.Warn(fmt.Sprintf("Reconnection failed : %v", token.Error()))
		} else {
			driver.Logger.Warn(fmt.Sprintf("Reconnection sucessful"))
		}
	})

	client := MQTT.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}

func newResult(deviceObject models.DeviceResource, ro models.ResourceOperation, reading interface{}) (*sdkModel.CommandValue, error) {
	var result = &sdkModel.CommandValue{}
	var err error
	var resTime = time.Now().UnixNano() / int64(time.Millisecond)
	var profileValueType = strings.ToLower(deviceObject.Properties.Value.Type)
	var readingValueType = strings.ToLower(reflect.TypeOf(reading).String())

	if readingValueType == "int" {
		reading = int64(reading.(int))
	}

	if readingValueType == "uint" {
		reading = uint64(reading.(uint))
	}

	// Check and convert reading when it is string type
	reading, readingValueType, err = handleReadingStringValue(profileValueType, readingValueType, reading)
	if err != nil {
		err = fmt.Errorf("parse reading fail. Error: %v", err)
		driver.Logger.Error(err.Error())
		return result, err
	}

	// Only check unknown value range
	if readingValueType == "int" || readingValueType == "uint" || readingValueType == "float32" || readingValueType == "float64" {
		if !checkValueInRange(profileValueType, readingValueType, reading) {
			err = fmt.Errorf("parse reading fail. Reading(%v) is out of the value type(%v)'s range", readingValueType, profileValueType)
			driver.Logger.Error(err.Error())
			return result, err
		}
	} else {
		// Throw error when value type no need to convert but not matched
		if profileValueType != readingValueType {
			err = fmt.Errorf("parse reading fail. Value type not matched. Reading - (%v) , profile - (%v)", readingValueType, profileValueType)
			driver.Logger.Error(err.Error())
			return result, err
		}
	}

	// Convert int, uint, float to correct value type
	reading = convertReadingValueType(profileValueType, readingValueType, reading)

	switch profileValueType {
	case "bool":
		result, err = sdkModel.NewBoolValue(&ro, resTime, reading.(bool))
	case "string":
		result = sdkModel.NewStringValue(&ro, resTime, reading.(string))
	case "uint8":
		result, err = sdkModel.NewUint8Value(&ro, resTime, reading.(uint8))
	case "uint16":
		result, err = sdkModel.NewUint16Value(&ro, resTime, reading.(uint16))
	case "uint32":
		result, err = sdkModel.NewUint32Value(&ro, resTime, reading.(uint32))
	case "uint64":
		result, err = sdkModel.NewUint64Value(&ro, resTime, reading.(uint64))
	case "int8":
		result, err = sdkModel.NewInt8Value(&ro, resTime, reading.(int8))
	case "int16":
		result, err = sdkModel.NewInt16Value(&ro, resTime, reading.(int16))
	case "int32":
		result, err = sdkModel.NewInt32Value(&ro, resTime, reading.(int32))
	case "int64":
		result, err = sdkModel.NewInt64Value(&ro, resTime, reading.(int64))
	case "float32":
		result, err = sdkModel.NewFloat32Value(&ro, resTime, reading.(float32))
	case "float64":
		result, err = sdkModel.NewFloat64Value(&ro, resTime, reading.(float64))
	default:
		err = fmt.Errorf("return result fail, none supported value type: %v", deviceObject.Properties.Value.Type)
	}

	return result, err
}

func newCommandValue(deviceObject models.DeviceResource, param *sdkModel.CommandValue) (interface{}, error) {
	var commandValue interface{}
	var err error
	switch strings.ToLower(deviceObject.Properties.Value.Type) {
	case "bool":
		commandValue, err = param.BoolValue()
	case "string":
		commandValue, err = param.StringValue()
	case "uint8":
		commandValue, err = param.Uint8Value()
	case "uint16":
		commandValue, err = param.Uint16Value()
	case "uint32":
		commandValue, err = param.Uint32Value()
	case "uint64":
		commandValue, err = param.Uint64Value()
	case "int8":
		commandValue, err = param.Int8Value()
	case "int16":
		commandValue, err = param.Int16Value()
	case "int32":
		commandValue, err = param.Int32Value()
	case "int64":
		commandValue, err = param.Int64Value()
	case "float32":
		commandValue, err = param.Float32Value()
	case "float64":
		commandValue, err = param.Float64Value()
	default:
		err = fmt.Errorf("fail to convert param, none supported value type: %v", deviceObject.Properties.Value.Type)
	}

	return commandValue, err
}

// fetchCommandResponse use to wait and fetch response from CommandResponses map
func fetchCommandResponse(commandResponses map[string]string, cmdUuid string) (string, bool) {
	var cmdResponse string
	var ok bool
	for i := 0; i < 5; i++ {
		cmdResponse, ok = commandResponses[cmdUuid]
		if ok {
			break
		} else {
			time.Sleep(time.Second * time.Duration(1))
		}
	}

	return cmdResponse, ok
}
