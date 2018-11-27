// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"encoding/json"
	"fmt"
	"net/url"
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

	switch deviceObject.Properties.Value.Type {
	case "Bool":
		result, err = sdkModel.NewBoolValue(&ro, resTime, reading.(bool))
	case "String":
		result = sdkModel.NewStringValue(&ro, resTime, reading.(string))
	case "Uint8":
		result, err = sdkModel.NewUint8Value(&ro, resTime, reading.(uint8))
	case "Uint16":
		result, err = sdkModel.NewUint16Value(&ro, resTime, reading.(uint16))
	case "Uint32":
		result, err = sdkModel.NewUint32Value(&ro, resTime, reading.(uint32))
	case "Uint64":
		result, err = sdkModel.NewUint64Value(&ro, resTime, reading.(uint64))
	case "Int8":
		result, err = sdkModel.NewInt8Value(&ro, resTime, reading.(int8))
	case "Int16":
		result, err = sdkModel.NewInt16Value(&ro, resTime, reading.(int16))
	case "Int32":
		result, err = sdkModel.NewInt32Value(&ro, resTime, reading.(int32))
	case "Int64":
		result, err = sdkModel.NewInt64Value(&ro, resTime, reading.(int64))
	case "Float32":
		result, err = sdkModel.NewFloat32Value(&ro, resTime, reading.(float32))
	case "Float64":
		result, err = sdkModel.NewFloat64Value(&ro, resTime, reading.(float64))
	default:
		err = fmt.Errorf("return result fail, none supported value type: %v", deviceObject.Properties.Value.Type)
	}

	return result, err
}

func newCommandValue(deviceObject models.DeviceObject, param *sdkModel.CommandValue) (interface{}, error) {
	var commandValue interface{}
	var err error
	switch deviceObject.Properties.Value.Type {
	case "Bool":
		commandValue, err = param.BoolValue()
	case "String":
		commandValue, err = param.StringValue()
	case "Uint8":
		commandValue, err = param.Uint8Value()
	case "Uint16":
		commandValue, err = param.Uint16Value()
	case "Uint32":
		commandValue, err = param.Uint32Value()
	case "Uint64":
		commandValue, err = param.Uint64Value()
	case "Int8":
		commandValue, err = param.Int8Value()
	case "Int16":
		commandValue, err = param.Int16Value()
	case "Int32":
		commandValue, err = param.Int32Value()
	case "Int64":
		commandValue, err = param.Int64Value()
	case "Float32":
		commandValue, err = param.Float32Value()
	case "Float64":
		commandValue, err = param.Float64Value()
	default:
		err = fmt.Errorf("return result fail, none supported value type: %v", deviceObject.Properties.Value.Type)
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
