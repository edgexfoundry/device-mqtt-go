// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"strings"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

func TestCreateConnectionInfo(t *testing.T) {
	schema := "tcp"
	host := "0.0.0.0"
	port := "1883"
	user := "admin"
	password := "password"
	clientId := "CommandPublisher"
	topic := "CommandTopic"
	protocols := map[string]models.ProtocolProperties{
		Protocol: {
			Schema:   schema,
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			ClientId: clientId,
			Topic:    topic,
		},
	}

	connectionInfo, err := CreateConnectionInfo(protocols)

	if err != nil {
		t.Fatalf("Fail to create connectionIfo. Error: %v", err)
	}
	if connectionInfo.Schema != schema || connectionInfo.Host != host || connectionInfo.Port != port ||
		connectionInfo.User != user || connectionInfo.Password != password || connectionInfo.ClientId != clientId ||
		connectionInfo.Topic != topic {
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
		IncomingSchema: "tcp", IncomingHost: "0.0.0.0", IncomingPort: "1883",
		IncomingUser: "admin", IncomingPassword: "public", IncomingQos: "0",
		IncomingKeepAlive: "3600", IncomingClientId: "IncomingDataSubscriber", IncomingTopic: "DataTopic",

		ResponseSchema: "tcp", ResponseHost: "0.0.0.0", ResponsePort: "1883",
		ResponseUser: "admin", ResponsePassword: "public", ResponseQos: "0",
		ResponseKeepAlive: "3600", ResponseClientId: "CommandResponseSubscriber", ResponseTopic: "ResponseTopic",
	}
	diverConfig, err := CreateDriverConfig(configs)
	if err != nil {
		t.Fatalf("Fail to load config, %v", err)
	}
	if diverConfig.IncomingSchema != configs[IncomingSchema] || diverConfig.IncomingHost != configs[IncomingHost] ||
		diverConfig.IncomingPort != 1883 || diverConfig.IncomingUser != configs[IncomingUser] ||
		diverConfig.IncomingPassword != configs[IncomingPassword] || diverConfig.IncomingQos != 0 ||
		diverConfig.IncomingKeepAlive != 3600 || diverConfig.IncomingClientId != configs[IncomingClientId] ||
		diverConfig.IncomingTopic != configs[IncomingTopic] ||
		diverConfig.ResponseSchema != configs[ResponseSchema] || diverConfig.ResponseHost != configs[ResponseHost] ||
		diverConfig.ResponsePort != 1883 || diverConfig.ResponseUser != configs[ResponseUser] ||
		diverConfig.ResponsePassword != configs[ResponsePassword] || diverConfig.ResponseQos != 0 ||
		diverConfig.ResponseKeepAlive != 3600 || diverConfig.ResponseClientId != configs[ResponseClientId] ||
		diverConfig.ResponseTopic != configs[ResponseTopic] {

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
