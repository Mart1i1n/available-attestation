#!/bin/bash
VAL_NUM=${VAL_NUM:-"16384"}

echo "export VAL_NUM=${VAL_NUM}" > .env.sh
source .env.sh
sed -i "s#TotalValidatorCount = 256#TotalValidatorCount = $VAL_NUM#g" ./modified_prysm/beacon-chain/forkchoice/doubly-linked-tree/allchain.go
docker build -t attacker:latest -f dockerfiles/attack.Dockerfile .
docker build -t geth:latest -f dockerfiles/geth.Dockerfile .
docker build -t strategy:latest -f dockerfiles/strategy.Dockerfile .
docker build -t txpress-pos:latest -f dockerfiles/txpress.Dockerfile .
docker build -t modified_beacon:latest -f dockerfiles/modified.beacon.Dockerfile .
docker build -t normal_beacon:latest -f dockerfiles/normal.beacon.Dockerfile .
docker build -t modified_validator:latest -f dockerfiles/modified.validator.Dockerfile .
docker build -t normal_validator:latest -f dockerfiles/normal.validator.Dockerfile .
