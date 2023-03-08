#!/bin/sh

docker rm -f lps01
docker image rm lps

source ./lps-env.sh up