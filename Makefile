#!/usr/bin/make -f
#export GO111MODULE = on
#export GOPROXY=https://goproxy.cn,direct
#export GOPRIVATE="gitlab.bianjie.ai"
#export GONOPROXY="gitlab.bianjie.ai"
#export GONOSUMDB="gitlab.bianjie.ai"

build: go.sum
ifeq ($(OS),Windows_NT)
	go build  -o build/open-api.exe .
else
	go build  -o build/open-api .
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go install  .

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local gitlab.bianjie.ai/avata/open-api