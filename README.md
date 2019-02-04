# qrss-squeak
Generate CW and DFCW tones for QRSS (very slow morse code) transmissions.

In the following examples the base frequency is the default of 1000Hz.
That is, use Upper Sideband mode (or USB DATA or USB PACKET if available)
and tune your transmitter 1000 Hz below the target frequency.
These examples have short dit times for demo; in QRSS you want to slow that down.

### Quick Examples

Divide 1.2 by the desired words-per-minute to get the "dit" time in seconds.

1.2 / 12 = 0.1, so for 12 WPM use -dit=0.1:

```
go run qrss.go -mode=cw -dit=0.1 vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

Repeating that every 20 seconds, at 0, 20, and 40 seconds after the full minute:
```
go run qrss.go -loop=20 --mode=cw --dit=0.1 vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

Repeating that every 20 seconds, at 5, 25, and 45 seconds after the full minute:
```
go run qrss.go loop_ofset=5 -loop=20 --mode=cw --dit=0.1 vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### Example with Dual Frequency:
This uses 100Hz of bandwidth so you can hear the different tones if you
listen to it.  You can tighten the bandwith like to 5hz ( -bw=5 ) for QRSS.
```
go run qrss.go -mode=df -bw=100 -dit=0.3 vvv de q1rss/b |
  paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

### Debug Info

Add `-x` to dump debug info and exit.  You can use this to get the duration of the signal.
In this eamples, it prints that it will last 17.7 seconds.

```
$ go run qrss.go -mode=df -bw=10 -dit=0.3 -x  vvv de q1rss/b
2019/02/02 15:05:38 Play: DFEmitter{text="vvv de q1rss/b",morse="...- ...- ...-   -.. .   --.- .---- .-. ... ... -..-. -... ",freq=0.0,width=10.0,dit=300ms,total=17.7s}
```

## Command Line Reference

### Usage

The basic usage is
```
go run qrss.go  -flags...  words...
```

The words are the text to transmit.  Multiple words are joined with spaces.
It may be just one call sign like "Q1RSS" or multiple words like "VVV DE Q1RSS/B".

The flags are written like "--mode=cw" or "-dit=30".  You may use either one or two
dashes at the front; it does not matter.  Any non-boolean flags require an argument,
usually using "=" to connect the value.  All flags have some default.  You can see
a list of all flags and their defaults by invoking the command with a bogus flag, like
```
go run qrss.go -help
```
For more details about the flag parse, try "godoc flag | less".

### Audio on Standard Output

The program outputs 16-bit mono-channel samples in little-endian format (the least significant
byte comes first).  The output is "raw" in the sense there is no header at the front.
This is commonly called "PCM" format, even though it has nothing at all to do with
pulse-coded modulation.  These can be consumed by the Pulse Audio player command
"paplay" with flags --channels=1 (it's monophonic, not stereo), --format=s16le
(the samples are 16-bit little-endian), and --raw (there is no header).

The flag --rate can be used to set the number of samples per second.
The default is --rate=44100 which is the best rate for many sound cards.
The "paplay" command also takes a --rate argument.

### Frequency and Bandwidth

The flag --freq sets the "base frequency" (in Hz) for the sounds being played.
By default this is 1000, so you can set your Upper Sideband transmitter dial
frequency 1000 Hz lower than you want the transmission to be located.

The flag --bw sets the bandwidth range of the output tones.  For instance, if you're
transmitting Dual Frequency CW, and you specify "-mode=df -freq=1500 -bw=8",
then the program will output tones at 1500 and 1508 Hz.  The simple CW mode
(-mode=cw) only uses tones at the base frequency (1500 Hz) and ignores the -bw value.
(The value specified with --bw just specifies the maximum range of tones we try
to produce, but your actuall on-the-air bandwidth will always be greater
because all changing transmissions produce some extra sideband splatter.)

### Looping

The program can either output the words once ("one shot mode") or repeat that
every so many seconds ("looping mode").  For one shot mode, do not specify
the -loop flag.

For looping mode, use -loop to specify how long the loop is in seconds.
For example, to repeat every 20 seconds, specify `-loop=20`,
or to repeat every 5 minutes, specify `-loop=300`.

If the loop time is a divisor of 60, it will repeat at the same times
within every full minute, starting at the full minute.   If the loop
time is a divisor of 3600, it will repeat at the same times within every
full hour, starting at the full hour.  You can add an offset to that,
so it starts N seconds after the full minute or hour by specifying
`-loop_offset=N`.  For instance, to transmit at the 5 and 35 second
mark within every full minute, use `-loop_offset=5`.

(Technically, we take the "UNIX time" in seconds modulo the `-loop`
value, and if that equals the `-loop_offset` value, the transmission will
start on that second.)

When looping, audio bytes are output to stdout for one transmission of
the words, and then the output stops until the next transmission.
The loop duration should be longer than the transmision duration
plus the overhead of the paplay latency, the operating system's sound device
buffering, and the unix pipe buffering.  If not, the synchronization
of the loop time can get messed up.

## Modes

### cw: Continuous Wave Morse Code.

Using `-mode=cw` will produce ordinary Continuous Wave Morse Code.
It's rate is determined by the -dit flag, specifying the duration of
one dit in seconds.  Divide 1.2 by the desired words per minute to
get the dit time.

### fs: Frequency Shift Morse Code.

Using `-mode=fs` will produce two-tone frequency-shift encoded Morse Code.
The higher tone (at --freq value plus --bw value) represents "mark" state,
and looks just like CW morse code.  The lower tone (at --freq value)
is the "space" state, and is transmitted during the gaps in the morse code.
Use --bw to specify the difference between those frequencies.

### df: Dual Frequency Morse Code.

Using `-mode=dt` will produce Morse Code in which the dahs are the same
duration as the dits, but dahs are at the higher frequency and dits are
at the base frequency.
Use --bw to specify the difference between those frequencies.

### tf: Triple Frequency Morse Code.

Using `-mode=tf` gives you a varient of Dual Freuquency in which the
gaps between letters or words is filled with a third frequency midway
between the other two.

### hell: Human-Readable Hellschreiber

Using `-mode=hell` produces human-readable characters drawn in the
radio spectrum waterfall or screen grabs.  It uses a 5x8 (that is, 5x7
plus descenders) raster character generator.  (This generator is currently not
really optimized for Hellschreiber modes, but comes from LCD graphics
definitions used in the open-source Arduboy game console.  We should
get better fonts sometime.)

Use `-bw` to specify the maximum distance of the 8 tones produced.
For instance, if `-bw=21`, the tones will be 0, 3, 6, 9,12, 15, 18,
and 21 Hz offset from the base -freq value.

## Internal Go Documentation

TODO: Document internals.

You can write new modes and combine existing modes in creative ways
using configurable emitters and mixers and cron components.
Making this easy to do was a goal of the project.
If you read qrss.go you will get an idea.  Look at the demo modes
( `-mode=demo-clock` and `-mode=demo-junk` ).
