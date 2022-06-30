#!/bin/bash -ex

EDGEX_STARTUP_DURATION=$(snapctl get startup-duration)

if [ -n "$EDGEX_STARTUP_DURATION" ]; then
  export EDGEX_STARTUP_DURATION
fi

EDGEX_STARTUP_INTERVAL=$(snapctl get startup-interval)

if [ -n "$EDGEX_STARTUP_INTERVAL" ]; then
  export EDGEX_STARTUP_INTERVAL
fi

# convert cmdline to string array
ARGV=($@)

# grab binary path
BINPATH="${ARGV[0]}"

# binary name == service name/key
SERVICE=$(basename "$BINPATH")
SERVICE_ENV="$SNAP_DATA/config/$SERVICE/res/$SERVICE.env"
TAG="edgex-$SERVICE."$(basename "$0")

if [ -f "$SERVICE_ENV" ]; then
    logger --tag=$TAG "sourcing $SERVICE_ENV"
    set -o allexport
    source "$SERVICE_ENV" set
    set +o allexport 
fi

exec "$@"
