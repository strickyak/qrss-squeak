#!/bin/bash

GOPATH=/home/pi/go go run qrss.go -freq=1000 -mode=fractal --tx_on=tx --tx_off="sleep 3; rx" | aplay --rate=44100 --channels=1 --format=S16_LE -t raw  /dev/stdin
