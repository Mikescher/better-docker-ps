#!/bin/bash

sed -i 's/const BETTER_DOCKER_PS_VERSION = ".*"/const BETTER_DOCKER_PS_VERSION = "'$(git describe --tags --abbrev=0 | sed "s/v//")'"/' "version.go"
