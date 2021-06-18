// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

// Constants related to protocol properties
const (
	Protocol     = "mqtt"
	CommandTopic = "CommandTopic"
)

// Constants related to MQTT security
const (
	AuthModeUsernamePassword = "usernamepassword"
	AuthModeNone             = "none"
)

// Constants related to custom configuration
const (
	CustomConfigSectionName = "MQTTBrokerInfo"
	WritableInfoSectionName = CustomConfigSectionName + "/Writable"
)
