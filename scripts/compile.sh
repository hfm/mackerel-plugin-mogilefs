#!/bin/bash
#
# Original from: http://deeeet.com/writing/2014/05/19/gox/

set -e

DIR=$(cd $(dirname ${0})/.. && pwd)
pushd ${DIR}

XC_ARCH=${XC_ARCH:-386 amd64}
XC_OS=${XC_OS:-linux}

COMMIT=$(git describe --tags --always)

rm -rf pkg/
gox \
    -ldflags "-X main.GitCommit \"${COMMIT}\"" \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"
popd
