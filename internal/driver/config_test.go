// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchCommandTopic(t *testing.T) {
	topic := "test-command-topic"

	var tests = []struct {
		name              string
		properties        map[string]models.ProtocolProperties
		expectedErrorKind errors.ErrKind
	}{
		{name: "valid", properties: map[string]models.ProtocolProperties{Protocol: {CommandTopic: topic}}, expectedErrorKind: ""},
		{name: "invalid, protocol properties is not defined", properties: map[string]models.ProtocolProperties{}, expectedErrorKind: errors.KindContractInvalid},
		{name: "invalid, property is not exist", properties: map[string]models.ProtocolProperties{Protocol: {}}, expectedErrorKind: errors.KindContractInvalid},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			commandTopic, err := fetchCommandTopic(testCase.properties)

			if testCase.expectedErrorKind != "" {
				require.Equal(t, errors.Kind(err), testCase.expectedErrorKind)
				return
			} else {
				require.NoError(t, err)
				assert.Equal(t, topic, commandTopic)
			}
		})
	}

}
