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

docker build -t wrangler-proxy . || exit 1
docker tag wrangler-proxy gscloudcz/wrangler-proxy:latest
docker run -it wrangler-proxy

yes_or_no "Push the image to Docker Hub?" && docker push gscloudcz/wrangler-proxy:latest
