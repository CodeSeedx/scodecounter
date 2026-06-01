#!/bin/bash
set -e
cd /root/.picoclaw/workspace/repo/scodecounter
go get github.com/spf13/cobra
go mod tidy
go build -o scodecounter .
echo "Build successful!"
ls -la scodecounter
