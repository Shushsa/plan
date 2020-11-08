#!/usr/bin/env bash
# Builds everything for both nix and windows platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/linux.64bit/pland ./cmd/pland
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/linux.64bit/plancli ./cmd/plancli

GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ./build/windows.64bit/pland.exe ./cmd/pland
GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ./build/windows.64bit/plancli.exe ./cmd/plancli
