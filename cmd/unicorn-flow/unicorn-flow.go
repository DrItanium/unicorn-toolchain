package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/DrItanium/neuron"
	"github.com/DrItanium/unicorn-toolchain/microcode"
	"github.com/DrItanium/unicorn-toolchain/sys"
	"github.com/DrItanium/unicornhat"
	"syscall"
	"time"
)

var hyperspeed = flag.Bool("hyperspeed", false, "Disable delay (still consumes delay bytes)")
var brightness = flag.Float64("brightness-factor", unicornhat.DefaultBrightness(), "Set brightness cap (0.0 - 1.0).\n\tWARNING: If you set this brightness too high you can cause retinal damage and I'm not responsible for that!!!")

var microcodeWord microcode.MicrocodeWord
var elements = make([]byte, 4)

func microsecond_delay(usec time.Duration) {
	if !*hyperspeed {
		time.Sleep(usec * time.Microsecond)
	}
}

type wideFunctionBody func(word *microcode.MicrocodeWord)

func shiftWord(input *bufio.Reader) {
	for i := 63; i > 0; i-- {
		microcodeWord[i].Pixel.R = microcodeWord[i-1].Pixel.R
		microcodeWord[i].Pixel.G = microcodeWord[i-1].Pixel.G
		microcodeWord[i].Pixel.B = microcodeWord[i-1].Pixel.B
		microcodeWord[i].Delay = microcodeWord[i-1].Delay
	}
	// now read a new pixel in
	err := loadPixel(input, &microcodeWord[0])
	if err != nil {
		neuron.StopRunning()
	}
}
func loadPixel(input *bufio.Reader, fld *microcode.MicrocodeField) error {
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		return err
	} else {
		fld.Pixel.R = elements[0]
		fld.Pixel.G = elements[1]
		fld.Pixel.B = elements[2]
		fld.Delay = elements[3]
		return nil
	}
}

func main() {
	defer unicornsys.Terminate(0)
	flag.Parse()
	neuron.StopRunningOnSignal(syscall.SIGINT)
	input := neuron.NewStandardInReader()
	if *brightness > 1.0 || *brightness < 0.0 {
		fmt.Println("Brightness out of range, using the default brightness")
		*brightness = unicornhat.DefaultBrightness()
	} else if *brightness > unicornhat.DefaultBrightness() {
		fmt.Println("WARNING: you've set the brightness higher than the default, this can get bright! Please don't look directly at it!")
	}
	unicornhat.Initialize(64, *brightness)
	unicornhat.ClearLEDBuffer()
	for neuron.IsRunning() {
		for x := 0; neuron.IsRunning() && x < 8; x++ {
			for y := 0; neuron.IsRunning() && y < 8; y++ {
				index := unicornhat.CoordinateToPosition(x, y)
				microcodeWord[index].UpdateNativePixel(index, microsecond_delay)
			}
		}
		unicornhat.Show()
		shiftWord(input)
	}
	for i := 63; i > 0; i-- {
		microcodeWord[i] = microcodeWord[i-1]
		microcodeWord[i].UpdateNativePixel(i, microsecond_delay)
	}
	microcodeWord[0].Pixel.R = 0
	microcodeWord[0].Pixel.G = 0
	microcodeWord[0].Pixel.B = 0
	microcodeWord[0].UpdateNativePixel(0, microsecond_delay)
	unicornhat.Show()
	for q := 0; q < 64; q++ {
		for i := 63; i > 0; i-- {
			microcodeWord[i] = microcodeWord[i-1]
			microcodeWord[i].UpdateNativePixel(i, microsecond_delay)
		}
		unicornhat.Show()
	}
}
