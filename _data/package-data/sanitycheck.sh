#!/bin/bash

cd "$(dirname "$0")"

version_tag="$(cd ../../ && git tag --sort=-v:refname | grep -P 'v[0-9\.]' | head -1 | cut -c2-)"

version_code="$(cd ../../ && cat consts/version.go | grep -oP 'BETTER_DOCKER_PS_VERSION = .*' | grep -oP '"[0-9\.]+"' | grep -oP '[0-9\.]+' )"


if [ "$version_tag" != "$version_code" ]; then

  echo "Git-Tag version and Code-const version do not match!"
  echo "[GIT-TAG] := $version_tag"
  echo "[GO-CODE] := $version_code"

  exit 1

else

  echo "Version ('$version_tag') ok"

fi