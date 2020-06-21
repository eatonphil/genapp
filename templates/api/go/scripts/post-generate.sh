#!/usr/bin/env sh

set -e
set -u
set -x
set -o pipefail

go get golang.org/x/tools/cmd/goimports
if [ ! -f ../go.mod ]; then
    (cd .. && go mod init {{ api.extra.repo }})
fi
goimports -w ./
