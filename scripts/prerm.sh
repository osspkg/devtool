#!/bin/bash


if [ -f "/etc/systemd/system/devtool.service" ]; then
    systemctl stop devtool
    systemctl disable devtool
    systemctl daemon-reload
fi
