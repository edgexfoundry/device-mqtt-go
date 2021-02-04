// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

const (
	Protocol = "mqtt"

	Schema          = "Schema"
	Host            = "Host"
	Port            = "Port"
	ClientId        = "ClientId"
	Topic           = "Topic"
	CredentialsPath = "CredentialsPath"

	// Driver config
	IncomingSchema          = "IncomingSchema"
	IncomingHost            = "IncomingHost"
	IncomingPort            = "IncomingPort"
	IncomingQos             = "IncomingQos"
	IncomingKeepAlive       = "IncomingKeepAlive"
	IncomingClientId        = "IncomingClientId"
	IncomingTopic           = "IncomingTopic"
	IncomingCredentialsPath = "IncomingCredentialsPath"

	ResponseSchema          = "ResponseSchema"
	ResponseHost            = "ResponseHost"
	ResponsePort            = "ResponsePort"
	ResponseQos             = "ResponseQos"
	ResponseKeepAlive       = "ResponseKeepAlive"
	ResponseClientId        = "ResponseClientId"
	ResponseTopic           = "ResponseTopic"
	ResponseCredentialsPath = "ResponseCredentialsPath"

	CredentialsRetryTime = "CredentialsRetryTime"
	CredentialsRetryWait = "CredentialsRetryWait"

	ConnEstablishingRetry = "ConnEstablishingRetry"
	ConnRetryWaitTime     = "ConnRetryWaitTime"
)
