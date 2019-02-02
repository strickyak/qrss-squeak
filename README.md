# qrss-squeak
Generate CW and DFCW tones for QRSS (very slow morse code) transmissions.

In the following examples the base frequency is 1000Hz.
That is, use USB mode and tune your radio 1000 Hz below the target frequency.
These examples have short dit times for demo; in QRSS you want to slow that down.

### Example with CW at 12 WPM:
```
go run qrss.go -mode=cw -dit=0.1 -ramp=0.02  vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

Repeating that every 20 seconds, at 0, 20, and 40 seconds after the full minute:
```
go run qrss.go -loop=20 --mode=cw --dit=0.1 --ramp=0.02  vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

Repeating that every 20 seconds, at 5, 25, and 45 seconds after the full minute:
```
go run qrss.go loop_ofset=5 -loop=20 --mode=cw --dit=0.1 --ramp=0.02  vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### Example with Dual Frequency:
```
go run qrss.go --mode=df --bw=10 --dit=0.3 --ramp=0.05   vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### Debug Info

Add --x to dump debug info and exit:

```
$ go run qrss.go --mode=df --bw=10 --dit=0.3 --ramp=0.05 --x  vvv de q1rss/b
2019/02/02 15:05:38 Play: DFEmitter{text="vvv de q1rss/b",morse="...- ...- ...-   -.. .   --.- .---- .-. ... ... -..-. -... ",freq=0.0,width=10.0,dit=300ms,total=17.7s}
```
