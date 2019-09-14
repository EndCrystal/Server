#!/bin/bash
set -ex
root="$(realpath .)"
mkdir -p test/plugins
cd test/plugins
function do_build() (
  target="$1".ecplugin
  source="$root"/plugin/"$1"
  go build -buildmode=plugin -o "$target" "$source"
  chmod a+x "$target"
)
do_build websocket
cd ..
go run ..
