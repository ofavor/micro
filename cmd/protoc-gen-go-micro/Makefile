### Makefile for micro
GO        := GO111MODULE=on go
GOBUILD   := GO111MODULE=on CGO_ENABLED=0 $(GO) build

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"


.PHONY: proto

proto:
	@cd ./example/proto/ && protoc --go_out=./ --micro_out=./ *.proto
