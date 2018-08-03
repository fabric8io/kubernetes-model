#!/bin/bash
set -x

# Exit on error
set -e

source cico_setup.sh

setup

# Build kubernetes-model image
docker build -t kubernetes-model .

CID=$(docker run --detach=true -t kubernetes-model)

build_run_tests