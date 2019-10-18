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
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go"
	sdkModel "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/spf13/cast"
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
	CommandResponses sync.Map
	Config           *configuration
}

func NewProtocolDriver() sdkModel.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
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
	var cmd = req.DeviceResourceName

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
	cmdResponse, ok := d.fetchCommandResponse(cmdUuid)
	if !ok {
		err = fmt.Errorf("can not fetch command response: method=%v cmd=%v", method, cmd)
		return result, err
	}

	driver.Logger.Info(fmt.Sprintf("Parse command response: %v", cmdResponse))

	var response map[string]interface{}
	json.Unmarshal([]byte(cmdResponse), &response)
	reading, ok := response[cmd]
	if !ok {
		err = fmt.Errorf("can not fetch command reading: method=%v cmd=%v", method, cmd)
		return result, err
	}

	result, err = newResult(req, reading)
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
	var cmd = req.DeviceResourceName

	data := make(map[string]interface{})
	data["uuid"] = cmdUuid
	data["method"] = method
	data["cmd"] = cmd

	commandValue, err := newCommandValue(req.Type, param)
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
	var cmdResponse interface{}
	var ok bool
	for i := 0; i < 5; i++ {
		cmdResponse, ok = d.CommandResponses.Load(cmdUuid)
		if ok {
			d.CommandResponses.Delete(cmdUuid)
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

func newResult(req sdkModel.CommandRequest, reading interface{}) (*sdkModel.CommandValue, error) {
	var result = &sdkModel.CommandValue{}
	var err error
	var resTime = time.Now().UnixNano()
	castError := "fail to parse %v reading, %v"

	if !checkValueInRange(req.Type, reading) {
		err = fmt.Errorf("parse reading fail. Reading %v is out of the value type(%v)'s range", reading, req.Type)
		driver.Logger.Error(err.Error())
		return result, err
	}

	switch req.Type {
	case sdkModel.Bool:
		val, err := cast.ToBoolE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewBoolValue(req.DeviceResourceName, resTime, val)
	case sdkModel.String:
		val, err := cast.ToStringE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result = sdkModel.NewStringValue(req.DeviceResourceName, resTime, val)
	case sdkModel.Uint8:
		val, err := cast.ToUint8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewUint8Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Uint16:
		val, err := cast.ToUint16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewUint16Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Uint32:
		val, err := cast.ToUint32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewUint32Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Uint64:
		val, err := cast.ToUint64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewUint64Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Int8:
		val, err := cast.ToInt8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewInt8Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Int16:
		val, err := cast.ToInt16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewInt16Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Int32:
		val, err := cast.ToInt32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewInt32Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Int64:
		val, err := cast.ToInt64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewInt64Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Float32:
		val, err := cast.ToFloat32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewFloat32Value(req.DeviceResourceName, resTime, val)
	case sdkModel.Float64:
		val, err := cast.ToFloat64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
		result, err = sdkModel.NewFloat64Value(req.DeviceResourceName, resTime, val)
	default:
		err = fmt.Errorf("return result fail, none supported value type: %v", req.Type)
	}

	return result, err
}

func newCommandValue(valueType sdkModel.ValueType, param *sdkModel.CommandValue) (interface{}, error) {
	var commandValue interface{}
	var err error
	switch valueType {
	case sdkModel.Bool:
		commandValue, err = param.BoolValue()
	case sdkModel.String:
		commandValue, err = param.StringValue()
	case sdkModel.Uint8:
		commandValue, err = param.Uint8Value()
	case sdkModel.Uint16:
		commandValue, err = param.Uint16Value()
	case sdkModel.Uint32:
		commandValue, err = param.Uint32Value()
	case sdkModel.Uint64:
		commandValue, err = param.Uint64Value()
	case sdkModel.Int8:
		commandValue, err = param.Int8Value()
	case sdkModel.Int16:
		commandValue, err = param.Int16Value()
	case sdkModel.Int32:
		commandValue, err = param.Int32Value()
	case sdkModel.Int64:
		commandValue, err = param.Int64Value()
	case sdkModel.Float32:
		commandValue, err = param.Float32Value()
	case sdkModel.Float64:
		commandValue, err = param.Float64Value()
	default:
		err = fmt.Errorf("fail to convert param, none supported value type: %v", valueType)
	}

	return commandValue, err
}

// fetchCommandResponse use to wait and fetch response from CommandResponses map
func (d *Driver) fetchCommandResponse(cmdUuid string) (string, bool) {
	var cmdResponse interface{}
	var ok bool
	for i := 0; i < 5; i++ {
		cmdResponse, ok = d.CommandResponses.Load(cmdUuid)
		if ok {
			d.CommandResponses.Delete(cmdUuid)
			break
		} else {
			time.Sleep(time.Second * time.Duration(1))
		}
	}

	return fmt.Sprintf("%v", cmdResponse), ok
}

func (d *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debug(fmt.Sprintf("Device %s is added", deviceName))
	return nil
}

func (d *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

func (d *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}
