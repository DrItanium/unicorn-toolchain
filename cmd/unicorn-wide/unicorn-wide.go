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

var fullpageUpdate = flag.Bool("fullpage", false, "Perform full page updates")
var columnByColumn = flag.Bool("column", false, "Perform column by column updates")
var rowByRow = flag.Bool("row", false, "Perform row by row updates")
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
func nativePixelUpdate(input *bufio.Reader, fullpage bool) {
	loadWordAndPerformAction(input, func(word *microcode.MicrocodeWord) {
		for i := 0; neuron.IsRunning() && i < 64; i++ {
			word[i].UpdateNativePixel(i, microsecond_delay)
			if !fullpage {
				unicornhat.Show()
			}
		}
		if fullpage {
			unicornhat.Show()
		}
	})
}
func lineByLineUpdate(input *bufio.Reader, doRowByRow bool) {
	loadWordAndPerformAction(input, func(word *microcode.MicrocodeWord) {
		for x := 0; neuron.IsRunning() && x < 8; x++ {
			for y := 0; neuron.IsRunning() && y < 8; y++ {
				var index int
				if doRowByRow {
					index = unicornhat.CoordinateToPosition(y, x)
				} else {
					index = unicornhat.CoordinateToPosition(x, y)
				}
				word[index].UpdateNativePixel(index, microsecond_delay)
			}
			unicornhat.Show()
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
		if *fullpageUpdate {
			nativePixelUpdate(input, true)
		} else if *columnByColumn {
			lineByLineUpdate(input, false)
		} else if *rowByRow {
			lineByLineUpdate(input, true)
		} else {
			nativePixelUpdate(input, false)
		}
	}
}
