// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ConnectionInfo struct {
	Schema   string
	Host     string
	Port     int
	User     string
	Password string
	ClientId string
	Topic    string
}

type configuration struct {
	IncomingSchema    string
	IncomingHost      string
	IncomingPort      int
	IncomingUser      string
	IncomingPassword  string
	IncomingQos       int
	IncomingKeepAlive int
	IncomingClientId  string
	IncomingTopic     string

	ResponseSchema    string
	ResponseHost      string
	ResponsePort      int
	ResponseUser      string
	ResponsePassword  string
	ResponseQos       int
	ResponseKeepAlive int
	ResponseClientId  string
	ResponseTopic     string

	ConnEstablishingRetry int
	ConnRetryWaitTime     int
}

// CreateDriverConfig use to load driver config for incoming listener and response listener
func CreateDriverConfig(configMap map[string]string) (*configuration, error) {
	// Set default value before loading new configuration
	config := &configuration{
		IncomingSchema:        DefaultSchema,
		IncomingHost:          DefaultHost,
		IncomingPort:          DefaultPort,
		IncomingUser:          DefaultUser,
		IncomingPassword:      DefaultPassword,
		IncomingQos:           DefaultQos,
		IncomingKeepAlive:     DefaultKeepAlive,
		IncomingClientId:      DefaultIncomingClientId,
		IncomingTopic:         DefaultIncomingTopic,
		ResponseSchema:        DefaultSchema,
		ResponseHost:          DefaultHost,
		ResponsePort:          DefaultPort,
		ResponseUser:          DefaultUser,
		ResponsePassword:      DefaultPassword,
		ResponseQos:           DefaultQos,
		ResponseKeepAlive:     DefaultKeepAlive,
		ResponseClientId:      DefaultResponseClientId,
		ResponseTopic:         DefaultResponseTopic,
		ConnEstablishingRetry: DefaultConnEstablishingRetry,
		ConnRetryWaitTime:     DefaultConnRetryWaitTime,
	}
	err := load(configMap, config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateConnectionInfo use to load MQTT connectionInfo for read and write command
func CreateConnectionInfo(protocols map[string]models.ProtocolProperties) (*ConnectionInfo, error) {
	// Set default value before loading new ConnectionInfo
	info := &ConnectionInfo{
		Schema:   DefaultSchema,
		Host:     DefaultHost,
		Port:     DefaultPort,
		User:     DefaultUser,
		Password: DefaultPassword,
		ClientId: DefaultCommandClientId,
		Topic:    DefaultCommandTopic,
	}
	protocol, ok := protocols[Protocol]
	if !ok {
		return info, fmt.Errorf("unable to load config, '%s' not exist", Protocol)
	}

	err := load(protocol, info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// load by reflect to check map key and then fetch the value
func load(config map[string]string, des interface{}) error {
	vals := reflect.ValueOf(des).Elem()
	for i := 0; i < vals.NumField(); i++ {
		typeField := vals.Type().Field(i)
		valueField := vals.Field(i)

		val, ok := config[typeField.Name]
		if !ok {
			continue
		}

		switch valueField.Kind() {
		case reflect.Int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			valueField.SetInt(int64(intVal))
		case reflect.String:
			valueField.SetString(val)
		default:
			return fmt.Errorf("none supported value type %v ,%v", valueField.Kind(), typeField.Name)
		}
	}
	return nil
}
