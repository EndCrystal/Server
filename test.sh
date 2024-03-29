#!/bin/bash
set -ex
root="$(realpath .)"
cd "$root"
mkdir -p test/plugins
cd test/plugins
function do_build() (
  target="$1".ecplugin
  source="$root"/plugin/"$1"
  go build -buildmode=plugin -o "$target" "$source"
  chmod a+x "$target"
)
do_build websocket
do_build flatdim
cd ..
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
(go run ../login) &
go run .. "$@"
