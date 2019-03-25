package marimba

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	. "github.com/strickyak/qrss-squeak/lib"
)

func loadImage(filename string) (image.Image, int, int) {
	r, err := os.Open(filename)
	Check(err)
	defer func() { r.Close() }()
	img, _, err := image.Decode(r)
	Check(err)
	b := img.Bounds()
	Must(b.Min.X == 0)
	Must(b.Min.Y == 0)
	return img, b.Max.X, b.Max.Y
}

type Conf struct {
	Dit       time.Duration
	Freq      float64
	Bandwidth float64
	Filename  string
	Gain      float64
}

type MarimbaEmitter struct {
	Conf
	Image image.Image
	W, H  int
}

func NewEmitter(conf *Conf) *MarimbaEmitter {
	im, w, h := loadImage(conf.Filename)
	o := &MarimbaEmitter{
		Conf:  *conf,
		Image: im,
		W:     w,
		H:     h,
	}
	return o
}

func (o *MarimbaEmitter) DurationInDits() float64 {
	return float64(o.W) + 2.0/3.0
}

func (o *MarimbaEmitter) Duration() time.Duration {
	return time.Duration(o.DurationInDits()) * o.Dit
}

func (o *MarimbaEmitter) DitPtr() *time.Duration {
	return &o.Dit
}

func (o *MarimbaEmitter) String() string {
	return fmt.Sprintf("MarimbaEmitter{%s(%dx%d),dit=%v,f=%v,bw=%v", o.Filename, o.W, o.H, o.Dit, o.Freq, o.Bandwidth)
}

func DurationTimesF(d time.Duration, n float64) time.Duration {
	return time.Duration(int64(float64(int64(d)) * n))
}
func DurationTimesI(d time.Duration, n int64) time.Duration {
	return time.Duration((int64(d)) * n)
}
func DurationTimesi(d time.Duration, n int) time.Duration {
	return time.Duration((int64(d)) * int64(n))
}

// Consumes 150% of dit duration:
// 0.25dit half-ramp-up  --- 0.25dit leader
// 0.25dit half-ramp-up  \
// 0.25dit full on        \
// 0.25dit full on         > one dit
// 0.25dit half-ramp-down /
// 0.25dit half-ramp-down --- 0.25dit trailer
func makeGap(numDits int, dit time.Duration, gapped bool) (Emitter, bool) {
	if gapped {
		numDits--
	}
	return NewCWEmitter(&CWConf{
		ToneWhenOff: false,
		Dit:         DurationTimesF(dit, float64(numDits)),
		Ramp:        0,
		Freq:        0,
		//Morse: []DiDah{' '},
		Text:  " ",
		Tail:  false,
		NoGap: true,
	}), false
}
func makeTone(numDits int, dit time.Duration, freq float64) (Emitter, bool) {
	e1 := NewCWEmitter(&CWConf{
		ToneWhenOff: false,
		Dit:         DurationTimesF(dit, 0.5+float64(numDits)),
		Ramp:        DurationTimesF(dit, 0.5),
		Freq:        freq,
		Text:        "e",
		Tail:        false,
		NoGap:       true,
	})
	e2 := NewCWEmitter(&CWConf{
		ToneWhenOff: false,
		Dit:         DurationTimesF(dit, 0.5),
		Ramp:        0,
		Freq:        freq,
		Text:        " ",
		Tail:        false,
		NoGap:       true,
	})
	return &Seq{
		Inputs: []Emitter{e1, e2},
	}, true
}

func pixel(x, y int, im image.Image) float64 {
	r, g, b, _ := im.At(x, y).RGBA()
	return float64(r+g+b) / (3.0 * 0xFFFF)
}

func (o *MarimbaEmitter) Emit(out chan Volt) {
	f1, f9 := o.Freq, o.Freq+o.Bandwidth
	bounds := o.Image.Bounds()
	x1, x9, y1, y9 := bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y
	Must(x1 == 0)
	Must(y1 == 0)

	var vec []Emitter
	for y := 0; y < y9; y++ {
		// for y := y9 - 1; y >= 0; y-- {
		seq := Seq{}
		//freq := f1 + (float64(y) * (f9 - f1) / float64(y9))
		// freq := f1 + (float64(y) * (f1 - f9) / float64(y9))
		freq := f1 + (float64(y9-y) * (f9 - f1) / float64(y9))

		gapped := false // Emitted 1 extra gap?  Not yet.
		prev := false
		count := 0
		for x := 0; x < x9; x++ {
			p := pixel(x, y, o.Image) > 0.5
			if p == prev {
				count++
			} else {
				var e Emitter
				if prev {
					e, gapped = makeTone(count, o.Dit, freq)
				} else {
					e, gapped = makeGap(count, o.Dit, gapped)
				}
				seq.Inputs = append(seq.Inputs, e)
				count = 1
			}
			prev = p
		}
		if prev {
			// Flush final tone.
			e, _ := makeTone(count, o.Dit, freq)
			seq.Inputs = append(seq.Inputs, e)
		}
		vec = append(vec, &seq)
	}
	mixer := &Mixer{Gain: o.Gain, Inputs: vec}
	log.Printf("MarimbaEmitter: %v", mixer)
	mixer.Emit(out)
}
