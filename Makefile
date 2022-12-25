.ONESHELL:

.DEFAULT_GOAL := all
.PHONY: install-protobuf check-os generate-protos

PROTOBUF_VERSION="3.14.0"
OS_NAME := $(shell uname -s)

ifeq ($(OS_NAME), Linux)
	OS="linux-x86_64"
else ifeq ($(OS_NAME), Darwin)
	OS="osx-x86_64"
endif

PROTOC_TARBALL="protoc-${PROTOBUF_VERSION}-${OS}.zip"
PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOC_TARBALL}"

all: check-os install-protobuf generate-protos
	@bash -c	"echo ${OS_NAME}"

install-protobuf:
	@mkdir -p tools
	if [ ! -d protobuf ]; then
		echo "OS is ${OS}"
		echo "Installing protoc ${PROTOBUF_VERSION} locally..."
		echo ${PROTOC_TARBALL}
		echo ${PROTOC_URL}

		pushd tools > /dev/null
		curl -LOs ${PROTOC_URL}
		unzip -qu ${PROTOC_TARBALL} -d protobuf
		rm ${PROTOC_TARBALL}

		echo "done."
		popd > /dev/null
	fi

check-os:
	@if [ ${OS_NAME} != "Linux" ] && [ ${OS_NAME} != "Darwin" ]; then
		@echo "Not recognized OS"
		@exit 1
	@fi

generate-protos:
	./generate_protos