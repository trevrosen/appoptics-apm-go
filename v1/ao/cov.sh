#!/bin/bash
set -e

COVERPKG="github.com/appoptics/appoptics-apm-go/v1/ao/internal/reporter,github.com/appoptics/appoptics-apm-go/v1/ao/internal/log,github.com/appoptics/appoptics-apm-go/v1/ao/internal/config,github.com/appoptics/appoptics-apm-go/v1/ao/internal/host,github.com/appoptics/appoptics-apm-go/v1/ao,github.com/appoptics/appoptics-apm-go/v1/ao/opentracing"
export APPOPTICS_DEBUG_LEVEL=1
go test -v -race -covermode=atomic -coverprofile=cov.out -coverpkg $COVERPKG

pushd internal/reporter/
go test -v -race -covermode=atomic -coverprofile=cov.out
popd

pushd internal/log/
go test -v -race -covermode=atomic -coverprofile=cov.out
popd

pushd internal/config/
go test -v -race -covermode=atomic -coverprofile=cov.out
popd

pushd internal/host/
go test -v -race -covermode=atomic -coverprofile=cov.out
popd

pushd opentracing
go test -v -race -covermode=atomic -coverprofile=cov.out
popd

gocovmerge cov.out internal/reporter/cov.out internal/log/cov.out internal/config/cov.out internal/host/cov.out opentracing/cov.out > covmerge.out

#go tool cover -html=covmerge.out
