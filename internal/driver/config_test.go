// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"strings"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/models"
)

func TestCreateConnectionInfo(t *testing.T) {
	schema := "tcp"
	host := "0.0.0.0"
	port := "1883"
	clientId := "CommandPublisher"
	topic := "CommandTopic"
	credentialsPath := "mqtt"

	protocols := map[string]models.ProtocolProperties{
		Protocol: {
			Schema:          schema,
			Host:            host,
			Port:            port,
			ClientId:        clientId,
			Topic:           topic,
			CredentialsPath: credentialsPath,
		},
	}

	connectionInfo, err := CreateConnectionInfo(protocols)

	if err != nil {
		t.Fatalf("Fail to create connectionIfo. Error: %v", err)
	}
	if connectionInfo.Schema != schema ||
		connectionInfo.Host != host ||
		connectionInfo.Port != port ||
		connectionInfo.ClientId != clientId ||
		connectionInfo.Topic != topic ||
		connectionInfo.CredentialsPath != credentialsPath {
		t.Fatalf("Unexpect test result. %v should match to %v ", connectionInfo, protocols)
	}
}

func TestCreateConnectionInfo_fail(t *testing.T) {
	protocols := map[string]models.ProtocolProperties{
		Protocol: {},
	}

	_, err := CreateConnectionInfo(protocols)
	if err == nil || !strings.Contains(err.Error(), "unable to load config") {
		t.Fatalf("Unexpect test result, config should be fail to load")
	}
}

func TestCreateDriverConfig(t *testing.T) {
	configs := map[string]string{
		IncomingSchema:          "tcp",
		IncomingHost:            "0.0.0.0",
		IncomingPort:            "1883",
		IncomingQos:             "0",
		IncomingKeepAlive:       "3600",
		IncomingClientId:        "IncomingDataSubscriber",
		IncomingTopic:           "DataTopic",
		IncomingCredentialsPath: "mqtt",
		ResponseSchema:          "tcp",
		ResponseHost:            "0.0.0.0",
		ResponsePort:            "1883",
		ResponseQos:             "0",
		ResponseKeepAlive:       "3600",
		ResponseClientId:        "CommandResponseSubscriber",
		ResponseTopic:           "ResponseTopic",
		ResponseCredentialsPath: "mqtt",
		CredentialsRetryTime:    "60",
		CredentialsRetryWait:    "1",
		ConnEstablishingRetry:   "10",
		ConnRetryWaitTime:       "5",
	}
	driverConfig, err := CreateDriverConfig(configs)
	if err != nil {
		t.Fatalf("Fail to load config, %v", err)
	}
	if driverConfig.IncomingSchema != configs[IncomingSchema] ||
		driverConfig.IncomingHost != configs[IncomingHost] ||
		driverConfig.IncomingPort != 1883 ||
		driverConfig.IncomingQos != 0 ||
		driverConfig.IncomingKeepAlive != 3600 ||
		driverConfig.IncomingClientId != configs[IncomingClientId] ||
		driverConfig.IncomingTopic != configs[IncomingTopic] ||
		driverConfig.IncomingCredentialsPath != configs[IncomingCredentialsPath] ||
		driverConfig.ResponseSchema != configs[ResponseSchema] ||
		driverConfig.ResponseHost != configs[ResponseHost] ||
		driverConfig.ResponsePort != 1883 ||
		driverConfig.ResponseQos != 0 ||
		driverConfig.ResponseKeepAlive != 3600 ||
		driverConfig.ResponseClientId != configs[ResponseClientId] ||
		driverConfig.ResponseTopic != configs[ResponseTopic] ||
		driverConfig.ResponseCredentialsPath != configs[ResponseCredentialsPath] ||
		driverConfig.CredentialsRetryTime != 60 ||
		driverConfig.CredentialsRetryWait != 1 ||
		driverConfig.ConnEstablishingRetry != 10 ||
		driverConfig.ConnRetryWaitTime != 5 {

		t.Fatalf("Unexpect test result, driver config doesn't correct load")
	}
}

func TestCreateDriverConfig_fail(t *testing.T) {
	configs := map[string]string{}
	_, err := CreateDriverConfig(configs)
	if err == nil || !strings.Contains(err.Error(), "unable to load config") {
		t.Fatalf("Unexpect test result, config should be fail to load")
	}
}
