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

	"github.com/edgexfoundry/device-sdk-go/pkg/service"

	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/config"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ConnectionInfo struct {
	Schema          string
	Host            string
	Port            string
	ClientId        string
	Topic           string
	CredentialsPath string
}

type configuration struct {
	IncomingSchema          string
	IncomingHost            string
	IncomingPort            int
	IncomingQos             int
	IncomingKeepAlive       int
	IncomingClientId        string
	IncomingTopic           string
	IncomingCredentialsPath string

	ResponseSchema          string
	ResponseHost            string
	ResponsePort            int
	ResponseQos             int
	ResponseKeepAlive       int
	ResponseClientId        string
	ResponseTopic           string
	ResponseCredentialsPath string

	CredentialsRetryTime int
	CredentialsRetryWait int

	ConnEstablishingRetry int
	ConnRetryWaitTime     int
}

// CreateDriverConfig use to load driver config for incoming listener and response listener
func CreateDriverConfig(configMap map[string]string) (*configuration, error) {
	config := new(configuration)
	err := load(configMap, config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateConnectionInfo use to load MQTT connectionInfo for read and write command
func CreateConnectionInfo(protocols map[string]models.ProtocolProperties) (*ConnectionInfo, error) {
	info := new(ConnectionInfo)
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
	errorMessage := "unable to load config, '%s' not exist"
	val := reflect.ValueOf(des).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)

		val, ok := config[typeField.Name]
		if !ok {
			return fmt.Errorf(errorMessage, typeField.Name)
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

func GetCredentials(secretPath string) (config.Credentials, error) {
	credentials := config.Credentials{}
	deviceService := service.RunningService()

	timer := startup.NewTimer(driver.Config.CredentialsRetryTime, driver.Config.CredentialsRetryWait)

	var secretData map[string]string
	var err error
	for timer.HasNotElapsed() {
		secretData, err = deviceService.SecretProvider.GetSecrets(secretPath, secret.UsernameKey, secret.PasswordKey)
		if err == nil {
			break
		}

		driver.Logger.Warnf(
			"Unable to retrieve MQTT credentials from SecretProvider at path '%s': %s. Retrying for %s",
			secretPath,
			err.Error(),
			timer.RemainingAsString())
		timer.SleepForInterval()
	}

	if err != nil {
		return credentials, err
	}

	credentials.Username = secretData[secret.UsernameKey]
	credentials.Password = secretData[secret.PasswordKey]

	return credentials, nil
}
