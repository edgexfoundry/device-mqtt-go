// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/config"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/models"
)

const CustomConfigSectionName = "MQTTBrokerInfo"

type ConnectionInfo struct {
	Schema          string
	Host            string
	Port            string
	ClientId        string
	Topic           string
	AuthMode        string
	CredentialsPath string
}

type ServiceConfig struct {
	MQTTBrokerInfo MQTTBrokerInfo
}

// UpdateFromRaw updates the service's full configuration from raw data received from
// the Service Provider.
func (sw *ServiceConfig) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ServiceConfig)
	if !ok {
		return false //errors.New("unable to cast raw config to type 'ServiceConfig'")
	}

	*sw = *configuration

	return true
}

type MQTTBrokerInfo struct {
	IncomingSchema          string
	IncomingHost            string
	IncomingPort            int
	IncomingQos             int
	IncomingKeepAlive       int
	IncomingClientId        string
	IncomingTopic           string
	IncomingAuthMode        string
	IncomingCredentialsPath string

	ResponseSchema          string
	ResponseHost            string
	ResponsePort            int
	ResponseQos             int
	ResponseKeepAlive       int
	ResponseClientId        string
	ResponseTopic           string
	ResponseAuthMode        string
	ResponseCredentialsPath string

	CredentialsRetryTime int
	CredentialsRetryWait int

	ConnEstablishingRetry int
	ConnRetryWaitTime     int

	// AuthMode is the MQTT broker authentication mechanism. Currently, 'usernamepassword' is the only AuthMode supported by this service, and the secret keys are 'username' and 'password'.
	AuthMode string
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

func SetCredentials(uri *url.URL, category string, authMode string, secretPath string) error {
	switch authMode {
	case AuthModeUsernamePassword:
		credentials, err := GetCredentials(secretPath)
		if err != nil {
			return fmt.Errorf("Unable to get %s MQTT credentials for secret path '%s': %s", category, secretPath, err.Error())
		}

		driver.Logger.Infof("%s MQTT credentials loaded", category)
		uri.User = url.UserPassword(credentials.Username, credentials.Password)

	case AuthModeNone:
		return nil
	default:
		return fmt.Errorf("invalid AuthMode '%s' for %s MQTT connection of", authMode, category)
	}

	return nil
}

func GetCredentials(secretPath string) (config.Credentials, error) {
	credentials := config.Credentials{}
	deviceService := service.RunningService()

	timer := startup.NewTimer(driver.serviceConfig.MQTTBrokerInfo.CredentialsRetryTime, driver.serviceConfig.MQTTBrokerInfo.CredentialsRetryWait)

	var secretData map[string]string
	var err error
	for timer.HasNotElapsed() {
		secretData, err = deviceService.SecretProvider.GetSecret(secretPath, secret.UsernameKey, secret.PasswordKey)
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
