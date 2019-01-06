#!/bin/bash
set -eu

function build() {
    echo "Building for OS=$1 ARCH=$2"
    env GOOS="$1" GOARCH="$2" go build -ldflags="-s -w" -o ${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-"$3"
    ! upx --ultra-brute ${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-"$3" || true
}

function install_deps() {
    sudo apt-get update
    sudo apt-get install -y upx dnsdist
}

function test_binary() {
    BINARY=${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-linux-amd64
    ${BINARY} -config examples/autoconf.toml -output /tmp/dnsdist.conf
    dnsdist  --supervised --disable-syslog -C /tmp/dnsdist.conf --check-config
}

function docker_build() {
    docker build . --tag jamesits/dnsdist-autoconf:azure-pipelines-latest
}

install_deps
build linux amd64 linux-amd64
# build windows amd64 windows-amd64.exe
test_binary
docker_build
