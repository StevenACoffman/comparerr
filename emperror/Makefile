SHELL := sh
.ONESHELL:
.EXPORT_ALL_VARIABLES:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

APP=emperrorerr
GOPATH?=${HOME}/go
INSTALLPATH=${GOPATH}/bin/${APP}

.PHONY: run
run: ## - Just run this
	go run main.go
.PHONY: build
build: ## - Just build it
	go build -o $(INSTALLPATH) -gcflags=-trimpath=$$HOME -asmflags=-trimpath=$$HOME main.go
	@echo Run $(INSTALLPATH)