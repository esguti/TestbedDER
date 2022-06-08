#/usr/bin/env bash

ps -eaf | grep modbus | head -1 | awk '{print $2}' | xargs -n1 kill
pkill arpspoof
