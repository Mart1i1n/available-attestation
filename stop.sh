#!/bin/bash
file=${1:-""}
docker compose -f $file down
