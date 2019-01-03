// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
	"github.com/edgexfoundry/edgex-go/pkg/models"
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

	config, err := LoadConfigFromFile()
	if err != nil {
		panic(fmt.Errorf("read MQTT driver configuration failed: %v", err))
	}
	d.Config = config

	go func() {
		err := startCommandResponseListening()
		if err != nil {
			panic(fmt.Errorf("start command response Listener failed: %v", err))
		}
	}()

	go func() {
		err := startIncomingListening()
		if err != nil {
			panic(fmt.Errorf("start incoming data Listener failed: %v", err))
		}
	}()

	return nil
}

func (d *Driver) DisconnectDevice(address *models.Addressable) error {
	panic("implement me")
}

func (d *Driver) HandleReadCommands(addr *models.Addressable, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	var responses = make([]*sdkModel.CommandValue, len(reqs))
	var err error

	// create device client and open connection
	var brokerUrl = addr.Address
	var brokerPort = addr.Port
	var username = addr.User
	var password = addr.Password
	var mqttClientId = addr.Publisher

	uri := &url.URL{
		Scheme: strings.ToLower(addr.Protocol),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createClient(mqttClientId, uri, 30)
	if err != nil {
		return responses, err
	}
	defer client.Disconnect(5000)

	for i, req := range reqs {
		res, err := d.handleReadCommandRequest(client, req, addr.Topic)
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
	var cmd = req.DeviceObject.Name

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
	reading, ok := response[req.DeviceObject.Name]
	if !ok {
		err = fmt.Errorf("can not fetch command reading: method=%v cmd=%v", method, cmd)
		return result, err
	}

	result, err = newResult(req.DeviceObject, req.RO, reading)
	if err != nil {
		return result, err
	} else {
		driver.Logger.Info(fmt.Sprintf("Get command finished: %v", result))
	}

	return result, err
}

func (d *Driver) HandleWriteCommands(addr *models.Addressable, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {
	var err error

	// create device client and open connection
	var brokerUrl = addr.Address
	var brokerPort = addr.Port
	var username = addr.User
	var password = addr.Password
	var mqttClientId = addr.Publisher

	uri := &url.URL{
		Scheme: strings.ToLower(addr.Protocol),
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createClient(mqttClientId, uri, 30)
	if err != nil {
		return err
	}
	defer client.Disconnect(5000)

	for i, req := range reqs {
		err = d.handleWriteCommandRequest(client, req, addr.Topic, params[i])
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
	var cmd = req.DeviceObject.Name

	data := make(map[string]interface{})
	data["uuid"] = cmdUuid
	data["method"] = method
	data["cmd"] = cmd

	commandValue, err := newCommandValue(req.DeviceObject, param)
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

func (*Driver) Stop(force bool) error {
	panic("implement me")
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

func newResult(deviceObject models.DeviceObject, ro models.ResourceOperation, reading interface{}) (*sdkModel.CommandValue, error) {
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

	// Check unknown value range
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

func newCommandValue(deviceObject models.DeviceObject, param *sdkModel.CommandValue) (interface{}, error) {
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

func checkValueInRange(profileValueType string, readingValueType string, reading interface{}) bool {
	isValid := false

	if profileValueType == "string" || profileValueType == "bool" {
		return true
	}

	if strings.Contains(profileValueType, "uint") {
		var val uint64
		if readingValueType == "int" {
			val = uint64(reading.(int64))
		} else if readingValueType == "uint" {
			val = uint64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = uint64(reading.(float32))
		} else {
			val = uint64(reading.(float64))
		}

		switch profileValueType {
		case "uint8":
			if val >= 0 && val <= math.MaxUint8 {
				isValid = true
			}
		case "uint16":
			if val >= 0 && val <= math.MaxUint16 {
				isValid = true
			}
		case "uint32":
			if val >= 0 && val <= math.MaxUint32 {
				isValid = true
			}
		case "uint64":
			maxiMum := uint64(math.MaxUint64)
			if val >= 0 && val <= maxiMum {
				isValid = true
			}
		}
		return isValid
	}

	if strings.Contains(profileValueType, "int") {
		var val int64
		if readingValueType == "int" {
			val = int64(reading.(int64))
		} else if readingValueType == "uint" {
			val = int64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = int64(reading.(float32))
		} else {
			val = int64(reading.(float64))
		}

		switch profileValueType {
		case "int8":
			if val >= math.MinInt8 && val <= math.MaxInt8 {
				isValid = true
			}
		case "int16":
			if val >= math.MinInt16 && val <= math.MaxInt16 {
				isValid = true
			}
		case "int32":
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				isValid = true
			}
		case "int64":
			if val >= math.MinInt64 && val <= math.MaxInt64 {
				isValid = true
			}
		}
		return isValid
	}

	if strings.Contains(profileValueType, "float") {
		var val float64
		if readingValueType == "int" {
			val = float64(reading.(int64))
		} else if readingValueType == "uint" {
			val = float64(reading.(uint64))
		} else if readingValueType == "float32" {
			val = float64(reading.(float32))
		} else {
			val = float64(reading.(float64))
		}

		switch profileValueType {
		case "float32":
			if val >= math.SmallestNonzeroFloat32 && val <= math.MaxFloat32 {
				isValid = true
			}
		case "float64":
			val := reading.(float64)
			if val >= math.SmallestNonzeroFloat64 && val <= math.MaxFloat64 {
				isValid = true
			}
		}

		return isValid
	}

	return isValid
}

func convertReadingValueType(profileValueType string, readingValueType string, reading interface{}) interface{} {
	if readingValueType == "int" {
		switch profileValueType {
		case "int8":
			reading = int8(reading.(int64))
		case "int16":
			reading = int16(reading.(int64))
		case "int32":
			reading = int32(reading.(int64))
		case "int64":
			reading = int64(reading.(int64))
		case "uint8":
			reading = uint8(reading.(int64))
		case "uint16":
			reading = uint16(reading.(int64))
		case "uint32":
			reading = uint32(reading.(int64))
		case "uint64":
			reading = uint64(reading.(int64))
		}
	} else if readingValueType == "uint" {
		switch profileValueType {
		case "int8":
			reading = int8(reading.(uint64))
		case "int16":
			reading = int16(reading.(uint64))
		case "int32":
			reading = int32(reading.(uint64))
		case "int64":
			reading = int64(reading.(uint64))
		case "uint8":
			reading = uint8(reading.(uint64))
		case "uint16":
			reading = uint16(reading.(uint64))
		case "uint32":
			reading = uint32(reading.(uint64))
		case "uint64":
			reading = uint64(reading.(uint64))
		}
	} else if readingValueType == "float64" {
		switch profileValueType {
		case "int8":
			reading = int8(reading.(float64))
		case "int16":
			reading = int16(reading.(float64))
		case "int32":
			reading = int32(reading.(float64))
		case "int64":
			reading = int64(reading.(float64))
		case "uint8":
			reading = uint8(reading.(float64))
		case "uint16":
			reading = uint16(reading.(float64))
		case "uint32":
			reading = uint32(reading.(float64))
		case "uint64":
			reading = uint64(reading.(float64))
		case "float32":
			reading = float32(reading.(float64))
		}
	} else if readingValueType == "float32" {
		switch profileValueType {
		case "int8":
			reading = int8(reading.(float32))
		case "int16":
			reading = int16(reading.(float32))
		case "int32":
			reading = int32(reading.(float32))
		case "int64":
			reading = int64(reading.(float32))
		case "uint8":
			reading = uint8(reading.(float32))
		case "uint16":
			reading = uint16(reading.(float32))
		case "uint32":
			reading = uint32(reading.(float32))
		case "uint64":
			reading = uint64(reading.(float32))
		case "float64":
			reading = float64(reading.(float32))
		}
	}

	return reading
}
