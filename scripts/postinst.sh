#!/bin/bash


if [ -f "/etc/systemd/system/devtool.service" ]; then
    systemctl start devtool
    systemctl enable devtool
    systemctl daemon-reload
fi
