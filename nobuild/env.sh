#!/bin/bash

# Copyright 2016 Huitse Tai. All rights reserved.
# Use of this source code is governed by BSD 3-clause
# license that can be found in the LICENSE file.

set -o errexit
set -o nounset
set -o pipefail

JORM_ROOT=$(dirname "${BASH_SOURCE}")/..

export GOPATH=`"${JORM_ROOT}/nobuild/gopath.sh"`
