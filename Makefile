GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)

makefile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
local_deploy_path := $(abspath $(dir $(makefile_path))/deploy)
base_dir := $(abspath $(shell git rev-parse --show-toplevel))

all:
	go build -o build/demo github.com/thirdmartini/gogui/cmd/demo
.PHONY: all

rpi:
	GOOS=linux GOARCH=arm64 go build -o build/demo.rpi github.com/thirdmartini/gogui/cmd/demo
.PHONY: rpi

run:
	go run github.com/thirdmartini/gogui/cmd/demo
.PHONY: run
