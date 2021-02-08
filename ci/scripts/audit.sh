#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-reporter-client
  nancy build && mv build/dp-reporter-client $cwd/build
popd