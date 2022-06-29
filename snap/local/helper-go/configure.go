// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package main

import (
	"fmt"
	"strings"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

// ConfToEnv defines mappings from snap config keys to EdgeX environment variable
// names that are used to override individual device-mqtt's [Driver]  configuration
// values via a .env file read by the snap service wrapper.
//
// The syntax to set a configuration key is:
//
// env.<section>.<keyname>
//
var ConfToEnv = map[string]string{
	// [MQTTBrokerInfo]
	"mqttbrokerinfo.schema":     "MQTTBROKERINFO_SCHEMA",
	"mqttbrokerinfo.host":       "MQTTBROKERINFO_HOST",
	"mqttbrokerinfo.port":       "MQTTBROKERINFO_PORT",
	"mqttbrokerinfo.qos":        "MQTTBROKERINFO_QOS",
	"mqttbrokerinfo.keep-alive": "MQTTBROKERINFO_KEEPALIVE",
	"mqttbrokerinfo.client-id":  "MQTTBROKERINFO_CLIENTID",

	"mqttbrokerinfo.credentials-retry-time":  "MQTTBROKERINFO_CREDENTIALSRETRYTIME",
	"mqttbrokerinfo.credentials-retry-wait":  "MQTTBROKERINFO_CREDENTIALSRETRYWAIT",
	"mqttbrokerinfo.conn-establishing-retry": "MQTTBROKERINFO_CONNESTABLISHINGRETRY",
	"mqttbrokerinfo.conn-retry-wait-time":    "MQTTBROKERINFO_CONNRETRYWAITTIME",

	"mqttbrokerinfo.auth-mode":        "MQTTBROKERINFO_AUTHMODE",
	"mqttbrokerinfo.credentials-path": "MQTTBROKERINFO_CREDENTIALSPATH",

	"mqttbrokerinfo.incoming-topic": "MQTTBROKERINFO_INCOMINGTOPIC",
	"mqttbrokerinfo.response-topic": "MQTTBROKERINFO_RESPONSETOPIC",

	// [Device]
	"device.update-last-connected": "DEVICE_UPDATELASTCONNECTED",
	"device.use-message-bus":       "DEVICE_USEMESSAGEBUS",
}

var cli *hooks.CtlCli = hooks.NewSnapCtl()

// configure is called by the main function
func configure() {
	var debug = false
	var err error
	var envJSON string

	status, err := cli.Config("debug")
	if err != nil {
		log.Fatalf("edgex-device-mqtt:configure: can't read value of 'debug': %v", err)
	}
	if status == "true" {
		debug = true
	}

	if err = hooks.Init(debug, "edgex-device-mqtt"); err != nil {
		log.Fatalf("edgex-device-mqtt:configure: initialization failure: %v", err)
	}

	// read env var override configuration
	envJSON, err = cli.Config(hooks.EnvConfig)
	if err != nil {
		log.Fatalf("Reading config 'env' failed: %v", err)
	}

	if envJSON != "" {
		hooks.Debug(fmt.Sprintf("edgex-device-mqtt:configure: envJSON: %s", envJSON))
		err = hooks.HandleEdgeXConfig("device-mqtt", envJSON, ConfToEnv)
		if err != nil {
			log.Fatalf("HandleEdgeXConfig failed: %v", err)
		}
	}

	// If autostart is not explicitly set, default to "no"
	// as only example service configuration and profiles
	// are provided by default.
	autostart, err := snapctl.Get("autostart").Run()
	if err != nil {
		log.Fatalf("Reading config 'autostart' failed: %v", err)
	}
	if autostart == "" {
		log.Debug("autostart is NOT set, initializing to 'no'")
		autostart = "no"
	}
	autostart = strings.ToLower(autostart)
	log.Debugf("autostart=%s", autostart)

	// services are stopped/disabled by default in the install hook
	switch autostart {
	case "true", "yes":
		err = snapctl.Start("device-mqtt").Enable().Run()
		if err != nil {
			log.Fatalf("Can't start service: %s", err)
		}
	case "false", "no":
		// no action necessary
	default:
		log.Fatalf("Invalid value for 'autostart': %s", autostart)
	}

	log.SetComponentName("configure")

	log.Info("Processing options")
	err = options.ProcessAppConfig("device-mqtt")
	if err != nil {
		log.Fatalf("could not process options: %v", err)
	}
}
