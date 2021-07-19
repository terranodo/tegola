#!/bin/sh -l

# This file expects to run in a Docker container running on 
# GitHub Actions. GitHub will mount automatically mount the source
# as a volumn inside this container with the following flg:
#
#	-v "/home/runner/work/tegola/tegola":"/github/workspace"
#
# The workdir is set using the following flag
#
#	--workdir /github/workspace
#
# The VERSION env var is automatically provied using the following flag:
# 
#	-e VERSION
#

# move to the tegola_lambda folder
cd cmd/tegola_lambda

# build the binary
go build -mod vendor -ldflags "-w -X main.Version=$VERSION"
