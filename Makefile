#!/usr/bin/make -f

build: go.sum
ifeq ($(OS),Windows_NT)
	go build  -o build/nftp.exe ./cmd/nftp
else
	go build  -o build/nftp ./cmd/nftp
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go install  ./cmd/nftp

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local gitlab.bianjie.ai/irita-paas/open-api