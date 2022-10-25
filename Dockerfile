# Copyright (c) 2020-2022 IOTech Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG BASE=golang:1.18-alpine3.16
FROM ${BASE} AS builder

ARG ADD_BUILD_TAGS=""
ARG MAKE="make -e ADD_BUILD_TAGS=$ADD_BUILD_TAGS build"
ARG ALPINE_PKG_BASE="make git openssh-client gcc libc-dev zeromq-dev libsodium-dev"
ARG ALPINE_PKG_EXTRA=""

# Install our build time packages.
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

WORKDIR /device-mqtt-go

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
# To run tests in the build container:
#   docker build --build-arg 'MAKE=build test' .
# This is handy of you do your Docker business on a Mac
RUN $MAKE

FROM alpine:3.16

LABEL license='VSPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2020-2021: IoTech Ltd'

# dumb-init needed for injected secure bootstrapping entrypoint script when run in secure mode.
RUN apk add --update --no-cache zeromq dumb-init

COPY --from=builder /device-mqtt-go/cmd /
COPY --from=builder /device-mqtt-go/LICENSE /
COPY --from=builder /device-mqtt-go/Attribution.txt /

EXPOSE 59982

ENTRYPOINT ["/device-mqtt"]
CMD ["--cp=consul://edgex-core-consul:8500", "--registry", "--confdir=/res"]
