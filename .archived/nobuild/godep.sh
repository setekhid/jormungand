#!/bin/bash

# Copyright 2016 Huitse Tai. All rights reserved.
# Use of this source code is governed by BSD 3-clause
# license that can be found in the LICENSE file.

set -o errexit
set -o nounset
set -o pipefail

JORM_ROOT=$(dirname "${BASH_SOURCE}")/..

JORM_DEPS=(
	"${JORM_ROOT}"
)

JORM_GOPATH_SRC="${JORM_ROOT}/.gopath/src"

jorm::nobuild::godep() {

	"${JORM_ROOT}/nobuild/tracedeps.go" ${JORM_DEPS[@]} | "${JORM_ROOT}/nobuild/copydeps.go" "${JORM_GOPATH_SRC}"
}

jorm::nobuild::godep
