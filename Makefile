.PHONY: test clean all

HOSTNAME=registry.terraform.io
NAMESPACE=shihanng
NAME=installer
VERSION=0.0.1
BINARY=terraform-provider-${NAME}
OS_ARCH ?= linux_amd64

build:
	goreleaser build --single-target --snapshot --rm-dist

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv dist/terraform-provider-${NAME}_${OS_ARCH}/* ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

test:
	go test $(TESTARGS) -race -parallel=4 ./...

testacc:
	TF_ACC=1 go test $(TESTARGS) -race -parallel=4 ./...
