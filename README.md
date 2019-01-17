# qrss-squeak
Generate CW and DFCW tones for QRSS (very slow morse code) transmissions.

In the following examples the base frequency is 1000Hz.
That is, use USB mode and tune your radio 1000 Hz below the target frequency.
These examples have short dit times for demo; in QRSS you want to slow that down.

### Example with CW at 12 WPM:
```
go run qrss.go --mode=cw --dit=0.1s --ramp=0.02   vvv de w6rek/b | \
  time paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### Example with Dual Frequency:
```
go run qrss.go --mode=df --step=10 --dit=0.3s --ramp=0.05   vvv de w6rek/b | \
  time paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### BUGS

*   Inconsistency where --dit is Duration but --ramp is float (seconds).
*   Duplicated code in mainCW and mainDF.
