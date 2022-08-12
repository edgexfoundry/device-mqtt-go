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
	hooks "github.com/canonical/edgex-snap-hooks/v2"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
)

// Deprecated
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

// configure is called by the main function
func configure() {
	const app = "device-mqtt"

	log.SetComponentName("configure")

	log.Info("Processing legacy env options")
	envJSON, err := hooks.NewSnapCtl().Config(hooks.EnvConfig)
	if err != nil {
		log.Fatalf("Reading config 'env' failed: %v", err)
	}
	if envJSON != "" {
		log.Debugf("envJSON: %s", envJSON)
		err = hooks.HandleEdgeXConfig(app, envJSON, ConfToEnv)
		if err != nil {
			log.Fatalf("HandleEdgeXConfig failed: %v", err)
		}
	}

	log.Info("Processing config options")
	err = options.ProcessConfig(app)
	if err != nil {
		log.Fatalf("could not process config options: %v", err)
	}

	log.Info("Processing autostart options")
	err = options.ProcessAutostart(app)
	if err != nil {
		log.Fatalf("could not process autostart options: %v", err)
	}
}
