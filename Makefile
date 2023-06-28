.PHONY: test clean all

NAME=installer
OS_ARCH ?= linux_amd64

build:
	goreleaser build --snapshot --clean

install: build
	rm -rf /tmp/tfproviders/
	mkdir -p /tmp/tfproviders/
	mv dist/terraform-provider-${NAME}_${OS_ARCH}/* /tmp/tfproviders/

test:
	go test $(TESTARGS) -race -parallel=4 ./...

testacc:
	TF_ACC=1 go test $(TESTARGS) -race -parallel=4 ./...
