#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
ethdir="$workspace/src/github.com/Skyhark-Projects"
if [ ! -L "$ethdir/golang-bic-from-iban" ]; then
    mkdir -p "$ethdir"
    cd "$ethdir"
    ln -s ../../../../../. golang-bic-from-iban
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$ethdir/golang-bic-from-iban"
PWD="$ethdir/golang-bic-from-iban"

# Launch the arguments with the configured environment.
exec "$@"
