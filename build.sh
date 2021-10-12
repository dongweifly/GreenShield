#!/usr/bin/env bash

PROG_NAME=$0
ACTION=$1

function usage() {
    echo "Usage: $PROG_NAME {build|clean|package}"
    exit 2
}

export GOPROXY=https://goproxy.io
export GO112MODULE=on
export GO111MODULE=on

APP_HOME="/opt/green_shield/"

PACKAGE="green_shield.tar.gz"

function exit_on_error() {
    exit_code=$1
    last_command=${@:2}
    if [[ $exit_code -ne 0 ]]; then
        >&2 echo "\"${last_command}\" command failed with exit code ${exit_code}."
        exit $exit_code
    fi
}

function clean() {
    go clean
    rm -f green_shield.tar.gz
    rm -f green_shield
    rm -rf build
}

function build() {
    go mod tidy
    go build -o green_shield
    exit_on_error $? !!
}

function package() {
    mkdir build
    cp green_shield build
    cp deploy.sh build
    if [ -d "conf" ]; then
        cp -r conf build
    fi

    if [ -d "sensitive-words" ]; then
        cp -r sensitive-words build
    fi

#  tar -zcvf ${PACKAGE} build/*
}

function install() {
    if [[ -d "${APP_HOME}" ]]; then
        mkdir -p ${APP_HOME}
    fi

    cp green_shield ${APP_HOME}
    cp deploy.sh ${APP_HOME}
    if [ -d "conf" ]; then
        cp -r conf ${APP_HOME}
    fi

    if [ -d "sensitive-words" ]; then
        cp -r sensitive-words ${APP_HOME}
    fi
}

case "$ACTION" in
    build)
        build
    ;;
    clean)
        clean
    ;;
    package)
        clean
        build
        package
    ;;
    install)
        clean
        build
        install
    ;;
    *)
      usage
    ;;
esac

