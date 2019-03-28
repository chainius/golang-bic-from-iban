#!/bin/bash

clear

docker build . -t skyhark/iban-validator:v$1 && \
docker push skyhark/iban-validator:v$1