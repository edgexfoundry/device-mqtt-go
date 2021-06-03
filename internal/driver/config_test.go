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
