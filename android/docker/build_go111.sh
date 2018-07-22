#!/bin/bash -l
set -e

eval "$(gimme)"
export GOROOT_BOOTSTRAP=$GOROOT
cd go/src
./make.bash
