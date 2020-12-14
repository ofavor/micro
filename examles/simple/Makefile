### Makefile for nano
GO        := GO111MODULE=on go
GOBUILD   := GO111MODULE=on CGO_ENABLED=0 $(GO) build

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"


.PHONY: proto

proto:
	@cd ./server/proto/ && protoc --go_out=./ *.proto
