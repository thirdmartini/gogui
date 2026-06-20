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

run-demo:
	go run github.com/thirdmartini/gogui/cmd/demo
.PHONY: run-demo

run:
	go run github.com/thirdmartini/gogui/cmd/composer-demo
.PHONY: run


runpi:
	go run github.com/thirdmartini/gogui/cmd/demo --driver=drm
.PHONY: runpi

vnc:
	/Applications/TigerVNC.app/Contents/MacOS/vncviewer localhost:9000 SendClipboard=false CursorType=system AlwaysCursor=true
.PHONY: vnc
