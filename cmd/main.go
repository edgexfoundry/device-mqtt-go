// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/v3/pkg/startup"

	device_mqtt "github.com/edgexfoundry/device-mqtt-go"
	"github.com/edgexfoundry/device-mqtt-go/internal/driver"
)

const (
	serviceName string = "device-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_mqtt.Version, sd)
}
