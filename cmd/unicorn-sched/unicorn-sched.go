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

func microsecond_delay(usec time.Duration) {
	if !*hyperspeed {
		time.Sleep(usec * time.Microsecond)
	}
}

type wideFunctionBody func(word *microcode.MicrocodeWord)

func loadWordAndPerformAction(input *bufio.Reader, body wideFunctionBody) {
	word, err := microcode.ReadMicrocodeWord(input)
	if err == nil {
		body(word)
	}
}
func lineByLineUpdate(input *bufio.Reader) {
	loadWordAndPerformAction(input, func(word *microcode.MicrocodeWord) {
		for x := 0; neuron.IsRunning() && x < 8; x++ {
			for y, y0 := 0, 4; neuron.IsRunning() && y0 < 8; y, y0 = y+1, y0+1 {
				index := unicornhat.CoordinateToPosition(x, y)
				index2 := unicornhat.CoordinateToPosition(x, y0)
				word[index].UpdateNativePixel(index, microsecond_delay)
				word[index2].UpdateNativePixel(index2, microsecond_delay)
				unicornhat.Show()
			}
		}
	})
}

func main() {
	defer unicornsys.Terminate(0)
	flag.Parse()
	neuron.StopRunningOnSignal(syscall.SIGINT)
	if *brightness > 1.0 || *brightness < 0.0 {
		fmt.Println("Brightness out of range, using the default brightness")
		*brightness = unicornhat.DefaultBrightness()
	} else if *brightness > unicornhat.DefaultBrightness() {
		fmt.Println("WARNING: you've set the brightness higher than the default, this can get bright! Please don't look directly at it!")
	}
	unicornhat.Initialize(64, *brightness)
	unicornhat.ClearLEDBuffer()
	input := neuron.NewStandardInReader()
	for neuron.IsRunning() {
		lineByLineUpdate(input)
	}
}
