#!/bin/bash

node socketio-client.js &

sleep 1

go run main.go >> gobmpd.log &
