#!/bin/bash
set -e
if [ -z "$1" ]
then
      echo "please pass release tag as 1nd argument"
      exit
fi
docker build  . -t ghcr.io/middleware-labs/log-patterns-miner:$1
docker push ghcr.io/middleware-labs/log-patterns-miner:$1