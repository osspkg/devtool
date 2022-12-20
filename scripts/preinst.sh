#!/bin/bash


if ! [ -d /var/lib/devtool/ ]; then
    mkdir /var/lib/devtool
fi

if [ -f "/etc/systemd/system/devtool.service" ]; then
    systemctl stop devtool
    systemctl disable devtool
    systemctl daemon-reload
fi
