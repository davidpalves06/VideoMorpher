#!/bin/bash

set -xe

go build cmd/godeoeffects/main.go

./main
