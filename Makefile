REPOSITORY	:= github.com/jaymccon/cloudctl

GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)
GOBUILD		:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build

default:	build

build:
		$(GOBUILD)

.PHONY: default build
