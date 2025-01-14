#!/bin/bash

set -xe

go build cmd/videomorpher/videomorpher.go

./main
