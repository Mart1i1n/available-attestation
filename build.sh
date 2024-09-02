#!/bin/bash
docker build -t attacker:latest -f dockerfile/attack.Dockerfile .
docker build -t geth:latest -f dockerfile/geth.Dockerfile .
docker build -t strategy:latest -f dockerfile/strategy.Dockerfile .
docker build -t modified_beacon:latest -f dockerfile/modified.beacon.Dockerfile .
docker build -t normal_beacon:latest -f dockerfile/normal.beacon.Dockerfile .
docker build -t modified_validator:latest -f dockerfile/modified.validator.Dockerfile .
docker build -t normal_validator:latest -f dockerfile/normal.validator.Dockerfile .
