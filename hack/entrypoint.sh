#!/bin/sh

echo "Distributing files"
if [ -d "/opt/cni/bin/" ] && [ -f "./kubedtnhack" ]; then
  cp ./kubedtnhack /opt/cni/bin/
fi

# sleep 1000days
sleep 1000d