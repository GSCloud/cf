#!/bin/bash

function yes_or_no () {
  while true
  do
    read -p "$* [y/N]: " yn
    case $yn in
      [Yy]*) return 0 ;;
      [Nn]*) return 1 ;;
      *)
      return 1 ;;
    esac
  done
}

docker build -t gscloudcz/wrangler-proxy:latest . || exit 1
docker run -it gscloudcz/wrangler-proxy:latest

yes_or_no "Push the image to Docker Hub?" && docker push gscloudcz/wrangler-proxy:latest
