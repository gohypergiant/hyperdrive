#!/bin/bash
curl -fsSL github.com/gohypergiant/hyperdrive/releases/latest/download/{hyper.zip} -O hyper.zip
unzip hyper.zip
chmod +x hyper
./hyper jupyter remoteHost