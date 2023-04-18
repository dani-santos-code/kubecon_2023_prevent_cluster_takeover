#!/bin/bash
set -euo pipefail

PROJECT=kubeaudit-demo
ZONE=europe-west1-c
PORTS=(31080)

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

gcloud config set project ${PROJECT}
gcloud config set compute/zone ${ZONE}

# install numpy on laptop to make IAP tunnel faster
$(gcloud info --format="value(basic.python_location)") -m pip install numpy

# proxy connections so we can connect to internal cluster service from our laptop
# disable connection check since some ports won't be listening yet when we run the command
INSTANCE_NAME=$(gcloud compute instances list --limit=1 --zones=${ZONE} --filter='name: gke*' --format='value(name)')
for port in "${PORTS[@]}"; do
  gcloud compute start-iap-tunnel \
    --iap-tunnel-disable-connection-check --zone=${ZONE} --local-host-port=localhost:${port} \
    ${INSTANCE_NAME} ${port} &
done

sleep 3600
