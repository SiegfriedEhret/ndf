#!/usr/bin/env bash
gox -output="bin/{{.OS}}/{{.Arch}}/{{.Dir}}" -osarch="linux/amd64"