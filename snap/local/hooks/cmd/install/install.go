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
	"os"
	"path/filepath"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
)

var cli *hooks.CtlCli = hooks.NewSnapCtl()

func installConfigAndProfiles() error {
	var err error

	path := "/config/device-mqtt/res/configuration.toml"
	destFile := hooks.SnapData + path
	srcFile := hooks.Snap + path
	dir := filepath.Dir(destFile)

	// if configuration.toml already exists, it's been
	// provided by a content interface, so no need to
	// make the directory, which would cause any files
	// provided by the content interface to be deleted.
	if _, err = os.Stat(destFile); err == nil {
		return nil
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err = hooks.CopyFile(srcFile, destFile); err != nil {
		return err
	}

	path = "/config/device-mqtt/res/profiles/mqtt.test.device.profile.yml"
	destFile = hooks.SnapData + path
	srcFile = hooks.Snap + path

	if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
		return err
	}

	if err = hooks.CopyFile(srcFile, destFile); err != nil {
		return err
	}

	path = "/config/device-mqtt/res/devices/mqtt.test.device.toml"
	destFile = hooks.SnapData + path
	srcFile = hooks.Snap + path

	if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
		return err
	}

	if err = hooks.CopyFile(srcFile, destFile); err != nil {
		return err
	}

	return nil
}

func main() {
	var err error

	if err = hooks.Init(false, "edgex-device-mqtt"); err != nil {
		fmt.Println(fmt.Sprintf("edgex-device-mqtt::install: initialization failure: %v", err))
		os.Exit(1)

	}

	err = installConfigAndProfiles()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-device-mqtt:install: %v", err))
		os.Exit(1)
	}

}
