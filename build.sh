#!/usr/bin/env bash

VERSION="0.3-$(git rev-parse --short HEAD)"
APPNAME="Kakigoori"
go build -ldflags "-X main.APPNAME \"$APPNAME\" -X main.APPVERSION \"$VERSION\""
