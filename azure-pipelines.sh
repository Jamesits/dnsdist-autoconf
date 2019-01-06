#!/bin/bash
set -eu

export DEBIAN_FRONTEND=noninteractive
export GOPATH=/tmp/go
export GOBIN=/tmp/go/bin

function build() {
    echo "Building for OS=$1 ARCH=$2"
    env GOOS="$1" GOARCH="$2" go build -ldflags="-s -w" -o ${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-"$3"
    ! upx --ultra-brute ${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-"$3" || true
}

function install_deps() {
    echo "Installing dependencies"
    # install master version of dnsdist
    curl https://repo.powerdns.com/FD380FBB-pub.asc | sudo apt-key add -
    sudo cp docker/dnsdist.perference /etc/apt/preferences.d/dnsdist
    sudo cp docker/pdns.list.xenial /etc/apt/sources.list.d/pdns.list

    sudo apt-get update
    sudo apt-get install -y upx dnsdist
}

function test_binary() {
    echo "Testing binary"
    BINARY=${BUILD_ARTIFACTSTAGINGDIRECTORY}/dnsdist-autoconf-linux-amd64
    ${BINARY} -config examples/autoconf.toml -output /tmp/dnsdist.conf
    dnsdist -V
    dnsdist --supervised --disable-syslog -C /tmp/dnsdist.conf --check-config
}

function docker_build() {
    echo "Testing docker build"
    docker build . --tag jamesits/dnsdist-autoconf:azure-pipelines-latest
}

install_deps
go get ./...
build linux amd64 linux-amd64
# build windows amd64 windows-amd64.exe
test_binary
docker_build
