.PHONY: test clean all

HOSTNAME=registry.terraform.io
NAMESPACE=shihanng
NAME=setupenv
VERSION=0.0.0
OS_ARCH=linux_amd64
BINARY=terraform-provider-${NAME}

build:
	goreleaser build --single-target --snapshot --rm-dist

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv dist/terraform-provider-${NAME}_${OS_ARCH}/* ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
