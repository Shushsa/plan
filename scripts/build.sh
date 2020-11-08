#!/usr/bin/env bash
# Builds all the required stuff
GO111MODULE=on go get
GO111MODULE=on go install ./cmd/pland
GO111MODULE=on go install ./cmd/plancli