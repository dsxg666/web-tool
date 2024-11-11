#!/bin/bash
pkill tool
rm -rf ./tool
go build -o tool main.go
nohup ./tool &