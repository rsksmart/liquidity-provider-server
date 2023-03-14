#!/bin/bash

clear_cache() {
    echo "Clearing cache for $1..."
    image_ids=$(docker images -q $1)
    if [[ -n "$image_ids" ]]; then
        docker rmi -f $image_ids
    else
        echo "No images found for $1"
    fi
}

clear_cache "lbc-deployer:latest"
clear_cache "lps:latest"

if [[ "$1" == "--all" ]]; then
    clear_cache "rskj:latest"
    clear_cache "bitcond:latest"
    clear_cache "mongo:4"
fi
