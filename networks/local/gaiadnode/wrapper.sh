#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/iris/${BINARY:-iris}
ID=${ID:-0}
LOG=${LOG:-iris.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'iris' E.g.: -e BINARY=iris_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export GAIADHOME="/iris/node${ID}/iris"

if [ -d "`dirname ${GAIADHOME}/${LOG}`" ]; then
  "$BINARY" --home "$GAIADHOME" "$@" | tee "${GAIADHOME}/${LOG}"
else
  "$BINARY" --home "$GAIADHOME" "$@"
fi

chmod 777 -R /iris

