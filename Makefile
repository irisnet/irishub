#!/usr/bin/make -f

include scripts/build/build.mk
include scripts/build/godoc.mk
include scripts/build/contract.mk
include scripts/build/protobuf.mk
include scripts/build/testing.mk
include scripts/build/linting.mk