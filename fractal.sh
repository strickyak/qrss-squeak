#!/bin/bash

ID=W6REK

function tx_on() {
	: rigctl
  echo tx >&2
}
function tx_off() {
	: rigctl
  sleep 3
  echo rx >&2
}
function play() {
	paplay --rate=44100 --channels=1 --format=s16le --raw  /dev/stdin
}

go run qrss.go -freq=200 -mode=fractal  | paplay --rate=44100 --channels=1 --format=s16le --raw  /dev/stdin

