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
	for x := 7; x >= 0; x-- {
		// select the previous coordinate
		oldX := x - 1
		for y := 7; y >= 0; y-- {
			this := unicornhat.CoordinateToPosition(x, y)
			var prev int
			update := false
			oldY := y - 1
			if oldY < 0 {
				// we've underflowed
				if oldX >= 0 {
					prev = unicornhat.CoordinateToPosition(oldX, 7)
					update = true
				}
			} else {
				prev = unicornhat.CoordinateToPosition(x, oldY)
				update = true
			}
			if update {
				microcodeWord[this].Pixel.R = microcodeWord[prev].Pixel.R
				microcodeWord[this].Pixel.G = microcodeWord[prev].Pixel.G
				microcodeWord[this].Pixel.B = microcodeWord[prev].Pixel.B
				microcodeWord[this].Delay = microcodeWord[prev].Delay
			}
		}
	}
	index := unicornhat.CoordinateToPosition(0, 0)
	if !neuron.IsRunning() {
		microcodeWord[index].Pixel.R = 0
		microcodeWord[index].Pixel.G = 0
		microcodeWord[index].Pixel.B = 0
		microcodeWord[index].Delay = 0
	} else {
		// now read a new pixel in
		err := loadPixel(input, &microcodeWord[index])
		if err != nil {
			neuron.StopRunning()
		}
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
func flush(input *bufio.Reader) {
	for i := 0; i < 64; i++ {
		shiftWord(input)
		for j := 0; j < 64; j++ {
			microcodeWord[j].UpdateNativePixel(j, microsecond_delay)
		}
		unicornhat.Show()
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
	flush(input)
}
