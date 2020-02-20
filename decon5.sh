#!/bin/bash

A=' /   /   //   ///   ////  /  / '
B=' /   /  /     /  /  /     / /  '
C=' / / /  ///   ///   ///   //   '
D=' // //  /  /  / /   /     / /  '
E=' /   /   //   /  /  ////  /  / '

go run qrss.go \
  -freq=1010 -mode=decon5 -dit=5 -bw=20 \
  -tx_on='rig T 1' -tx_off='rig T 0' \
  -- "$E,$D,$C,$B,$A" |
	paplay --rate=44100 --channels=1 --format=s16le --raw  /dev/stdin

# 40m:  7038.750 + 1.010 kc
# 30m: 10139.000 + 1.010 kc
