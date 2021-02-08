#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-reporter-client
  make build && mv build/dp-reporter-client $cwd/build
  cp Dockerfile.concourse $cwd/build
popd