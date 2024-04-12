#!/bin/bash

arg1="$1"
arg2="$2"

echo "Starting supervisord..."
supervisord -c /etc/supervisord.conf -s

sleep 3

echo "Running sidecar-shim..."
exec /sidecar-shim $arg1 $arg2