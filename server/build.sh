#!/bin/bash
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
go build -ldflags="-s -w -X 'main.Version=$VERSION'" -o TRwindServer