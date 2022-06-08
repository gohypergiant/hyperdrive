#!/bin/bash
curl -fsSL https://github.com/gohypergiant/hyperdrive/archive/refs/tags/v0.0.0.zip -o hyper.zip

unzip hyper.zip -d /tmp/hyperdrive
cd /tmp/hyperdrive/hyperdrive-0.0.0/
chmod +x /tmp/hyperdrive/hyperdrive-0.0.0/hyper/hyper
mv /tmp/hyperdrive/hyperdrive-0.0.0/hyper/hyper /usr/bin/hyper
hyper jupyter remoteHost
