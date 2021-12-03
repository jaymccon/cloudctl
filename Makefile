REPOSITORY	:= github.com/jaymccon/cloudctl

GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)
GOBUILD		:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build

PROFILE     := default
REGION      := $(shell aws configure get region --profile $(PROFILE) || echo us-east-1)

default:	build

build:
		$(GOBUILD)

get-schemas:
		mkdir -p schemas
		$(foreach t, $(shell aws cloudformation list-types --type RESOURCE --visibility PUBLIC --provisioning-type FULLY_MUTABLE --filters Category=AWS_TYPES --query TypeSummaries[].TypeArn --output text --region $(REGION) --profile $(PROFILE)), $(shell aws cloudformation describe-type --arn $(t) --region $(REGION) --profile $(PROFILE) --query Schema --output text > schemas/$(word 3,$(subst /, ,$(t))).json && json-dereference -s schemas/$(word 3,$(subst /, ,$(t))).json -o schemas/$(word 3,$(subst /, ,$(t))).json))

.PHONY: default build
