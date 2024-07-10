#!/bin/bash

# test temp and temp/ubuntu-22.04.qcow2
mkdir -p temp

#  Download the Ubuntu 22.04 cloud image
if [ ! -f temp/ubuntu-22.04.qcow2 ]; then
  wget -O temp/ubuntu-22.04.qcow2 https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img
fi

kubectl virt image-upload pvc ubuntu2204-test \
--image-path=temp/ubuntu-22.04.qcow2 \
--size=50G \
--uploadproxy-url=https://10.0.0.31:31001 \
--insecure \
--wait-secs=60
 