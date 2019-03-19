// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
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
