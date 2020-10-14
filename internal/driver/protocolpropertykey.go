// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

const (
	Protocol = "mqtt"

	Schema   = "Schema"
	Host     = "Host"
	Port     = "Port"
	User     = "User"
	Password = "Password"
	ClientId = "ClientId"
	Topic    = "Topic"

	// Driver config
	IncomingSchema    = "IncomingSchema"
	IncomingHost      = "IncomingHost"
	IncomingPort      = "IncomingPort"
	IncomingUser      = "IncomingUser"
	IncomingPassword  = "IncomingPassword"
	IncomingQos       = "IncomingQos"
	IncomingKeepAlive = "IncomingKeepAlive"
	IncomingClientId  = "IncomingClientId"
	IncomingTopic     = "IncomingTopic"

	ResponseSchema    = "ResponseSchema"
	ResponseHost      = "ResponseHost"
	ResponsePort      = "ResponsePort"
	ResponseUser      = "ResponseUser"
	ResponsePassword  = "ResponsePassword"
	ResponseQos       = "ResponseQos"
	ResponseKeepAlive = "ResponseKeepAlive"
	ResponseClientId  = "ResponseClientId"
	ResponseTopic     = "ResponseTopic"

	ConnEstablishingRetry = "ConnEstablishingRetry"
	ConnRetryWaitTime     = "ConnRetryWaitTime"

	DefaultSchema                = "tcp"
	DefaultHost                  = "0.0.0.0"
	DefaultPort                  = 1883
	DefaultUser                  = ""
	DefaultPassword              = ""
	DefaultQos                   = 0
	DefaultKeepAlive             = 3600
	DefaultIncomingClientId      = "IncomingDataSubscriber"
	DefaultIncomingTopic         = "DataTopic"
	DefaultResponseClientId      = "CommandResponseSubscriber"
	DefaultResponseTopic         = "ResponseTopic"
	DefaultCommandClientId       = "CommandPublisher"
	DefaultCommandTopic          = "CommandTopic"
	DefaultConnEstablishingRetry = 10
	DefaultConnRetryWaitTime     = 5
)
