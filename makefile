INFO = github.com/mmbros/taskengine-app/cmd

APP_VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)

TEPKG_VERSION := $(shell grep "taskengine " go.mod | cut -d' ' -f2)

# GOXVER := $(shell go version | awk '{print $$3}')
GO_VERSION := $(shell go version)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date '+%F %T %z')
OS_ARCH := $(shell uname -s -m)

TIMESTAMP := $(shell date +%Y%m%dT%H%M%S)
# go build -v -ldflags="-X 'main.Version=v1.0.0' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'"


COMMON_LDFLAGS = -X '${INFO}.BuildTime=${BUILD_TIME}' \
                 -X '${INFO}.GitCommit=${GIT_COMMIT}' \
				 -X '${INFO}.GoVersion=${GO_VERSION}' \
				 -X '${INFO}.OsArch=${OS_ARCH}' \
				 -X '${INFO}.TaskEngineVersion=${TEPKG_VERSION}'

PROD_LDFLAGS = -ldflags "-X '${INFO}.AppVersion=${APP_VERSION}' ${COMMON_LDFLAGS}"

DEV_LDFLAGS = -ldflags "-X '${INFO}.AppVersion=dev-${TIMESTAMP}' ${COMMON_LDFLAGS}"

BIN=taskengine-app

all: prod

dev:
	go build ${DEV_LDFLAGS} -o ${BIN} *.go

prod:
	go build ${PROD_LDFLAGS} -o ${BIN} *.go
