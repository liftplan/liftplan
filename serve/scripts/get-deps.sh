#!/bin/bash

set -e


OUTPUT_DIR=static
PICO_VERSION=v2.0.6
PICO_SPEC=pico.fluid.classless.conditional.zinc.min
FIXI_VERSION=0.6.4

curl "https://github.com/picocss/pico/blob/${PICO_VERSION}/css/${PICO_SPEC}.css" > "${OUTPUT_DIR}/${PICO_SPEC}.${PICO_VERSION}.css"
curl "https://raw.githubusercontent.com/bigskysoftware/fixi/refs/tags/${FIXI_VERSION}/fixi.js" > "${OUTPUT_DIR}/fixi-${FIXI_VERSION}.js"
#openssl sha256 -binary < "${OUTPUT_DIR}/fixi-${FIXI_VERSION}.js" | openssl base64 > "${OUTPUT_DIR}/fixi-integrity.txt"
