#!/bin/bash
# Grab dependencies.

set -eu

# Ensure we only have one entry in GOPATH (glock gets confused
# by more).
export GOPATH=$(cd $(dirname $0)/../../../../../.. && pwd)

# go vet is special: it installs into $GOROOT (which $USER may not have
# write access to) instead of $GOPATH. It is usually but not always
# installed along with the rest of the go toolchain. Don't try to
# install it if it's already there.
if ! go vet 2>/dev/null; then
    go get golang.org/x/tools/cmd/vet
fi

if ! test -e ${GOPATH}/bin/glock ; then
    # glock is used to manage the rest of our dependencies (and to update
    # itself, so no -u here)
    go get github.com/robfig/glock
fi

${GOPATH}/bin/glock sync github.com/cockroachdb/cockroach

set -x

(cd ui && npm install)
