# Device MQTT Go
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-mqtt-go/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-mqtt-go/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/device-mqtt-go/branch/main/graph/badge.svg?token=IUywg34zfH)](https://codecov.io/gh/edgexfoundry/device-mqtt-go) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-mqtt-go)](https://goreportcard.com/report/github.com/edgexfoundry/device-mqtt-go) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-mqtt-go?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-mqtt-go/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-mqtt-go?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-mqtt-go)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-mqtt-go) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-mqtt-go)](https://github.com/edgexfoundry/device-mqtt-go/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-mqtt-go)](https://github.com/edgexfoundry/device-mqtt-go/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-mqtt-go-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-mqtt-go)](https://github.com/edgexfoundry/device-mqtt-go/commits)

## Overview
MQTT Micro Service - device service for connecting a MQTT topic to EdgeX acting like a device/sensor feed.

## Build with NATS Messaging
Currently, the NATS Messaging capability (NATS MessageBus) is opt-in at build time.
This means that the published Docker image and Snaps do not include the NATS messaging capability.

The following make commands will build the local binary or local Docker image with NATS messaging
capability included.
```makefile
make build-nats
make docker-nats
```

The locally built Docker image can then be used in place of the published Docker image in your compose file.
See [Compose Builder](https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder#gen) `nat-bus` option to generate compose file for NATS and local dev images.

## Packaging
This component is packaged as docker image and snap.

For docker, please refer to the [Dockerfile](Dockerfile) and [Docker Compose Builder](https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder) scripts.

For the snap, refer to the [snap](snap) directory.

## Usage
Users can refer to [this document](https://docs.edgexfoundry.org/2.1/examples/Ch-ExamplesAddingMQTTDevice) to learn how to use this device service.

## Community
- Discussion: https://github.com/orgs/edgexfoundry/discussions
- Mailing lists: https://lists.edgexfoundry.org/mailman/listinfo

## License
[Apache-2.0](LICENSE)
