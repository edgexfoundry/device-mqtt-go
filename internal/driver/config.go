// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type configuration struct {
	Incoming SubscribeInfo
	Response SubscribeInfo
}
type SubscribeInfo struct {
	Protocol     string
	Host         string
	Port         int
	Username     string
	Password     string
	Qos          int
	KeepAlive    int
	MqttClientId string
	Topic        string
}

type ConnectionInfo struct {
	Schema   string
	Host     string
	Port     string
	User     string
	Password string
	ClientId string
	Topic    string
}

// LoadConfigFromFile use to load toml configuration
func LoadConfigFromFile() (*configuration, error) {
	config := new(configuration)

	confDir := flag.Lookup("confdir").Value.(flag.Getter).Get().(string)
	if len(confDir) == 0 {
		confDir = flag.Lookup("c").Value.(flag.Getter).Get().(string)
	}

	if len(confDir) == 0 {
		confDir = "./res"
	}

	filePath := fmt.Sprintf("%v/configuration-driver.toml", confDir)

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("could not load configuration file (%s): %v", filePath, err.Error())
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return config, fmt.Errorf("unable to parse configuration file (%s): %v", filePath, err.Error())
	}
	return config, err
}

func CreateConnectionInfo(protocols map[string]models.ProtocolProperties) (info ConnectionInfo, err error) {
	errorMessage := "unable to create connection info, protocol config '%s' not exist"
	protocol, ok := protocols[Protocol]
	if !ok {
		return info, fmt.Errorf(errorMessage, Protocol)
	}
	schema, ok := protocol[Schema]
	if !ok {
		return info, fmt.Errorf(errorMessage, Schema)
	}
	host, ok := protocol[Host]
	if !ok {
		return info, fmt.Errorf(errorMessage, Host)
	}
	port, ok := protocol[Port]
	if !ok {
		return info, fmt.Errorf(errorMessage, Port)
	}
	user, ok := protocol[User]
	if !ok {
		return info, fmt.Errorf(errorMessage, User)
	}
	password, ok := protocol[Password]
	if !ok {
		return info, fmt.Errorf(errorMessage, Password)
	}
	clientId, ok := protocol[ClientId]
	if !ok {
		return info, fmt.Errorf(errorMessage, ClientId)
	}
	topic, ok := protocol[Topic]
	if !ok {
		return info, fmt.Errorf(errorMessage, Topic)
	}

	info = ConnectionInfo{
		Schema: schema, Host: host, Port: port,
		User: user, Password: password,
		ClientId: clientId, Topic: topic,
	}
	return info, nil
}
