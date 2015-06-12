#!/bin/bash
gox -osarch="linux/386 darwin/amd64" -output="bin/{{.Dir}}_{{.OS}}"