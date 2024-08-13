#!/bin/bash

TOKEN=$1
sed -i "s/GITHUB_TOKEN=/GITHUB_TOKEN=$TOKEN/g" ../../sample-config.env
sed -i 's/ENABLE_MANAGEMENT_API=false/ENABLE_MANAGEMENT_API=true/g' ../../sample-config.env