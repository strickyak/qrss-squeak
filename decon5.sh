#!/bin/bash

A=' /   /   //   ///   ////  /  / '
B=' /   /  /     /  /  /     / /  '
C=' / / /  ///   ///   ///   //   '
D=' // //  /  /  / /   /     / /  '
E=' /   /   //   /  /  ////  /  / '

go run qrss.go \
  -freq=1010 -mode=decon5 -dit=5 -bw=20 \
  -tx_on='ic-7100 T 1' -tx_off='ic-7100 T 0' \
  -- "$E,$D,$C,$B,$A" |
	paplay --rate=44100 --channels=1 --format=s16le --raw  /dev/stdin


# 40m:  7.038.750 + 1010
# 30m: 10.139.000 + 1010

#go run qrss.go \
#  -freq=1000 -mode=decon5 -dit=5 -bw=40 \
#  -tx_on='ic-7100 T 1' -tx_off='ic-7100 T 0' \
#  -- "$E,$D,$C,$B,$A" |
#	paplay --rate=44100 --channels=1 --format=s16le --raw  /dev/stdin
