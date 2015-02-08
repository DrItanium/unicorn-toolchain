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
	if err != nil {
		return
	}
	body(word)
}
func nativePixelUpdate(input *bufio.Reader, fullpage bool) {
	loadWordAndPerformAction(input, func(word *microcode.MicrocodeWord) {
		for i := 0; i < 64; i++ {
			if !neuron.IsRunning() {
				return
			}
			unicornhat.SetPixelColorType(uint(i), &word[i].Pixel)
			if !fullpage {
				unicornhat.Show()
			}
			microsecond_delay(time.Duration(word[i].Delay))
		}
		if fullpage {
			unicornhat.Show()
		}
	})
}
func lineByLineUpdate(input *bufio.Reader, doRowByRow bool) {
	loadWordAndPerformAction(input, func(word *microcode.MicrocodeWord) {
		for x := 0; x < 8; x++ {
			if !neuron.IsRunning() {
				return
			}
			for y := 0; y < 8; y++ {
				if !neuron.IsRunning() {
					return
				}
				var index int
				if doRowByRow {
					index = unicornhat.CoordinateToPosition(y, x)
				} else {
					index = unicornhat.CoordinateToPosition(x, y)
				}
				unicornhat.SetPixelColorType(uint(index), &word[index].Pixel)
				microsecond_delay(time.Duration(word[index].Delay))
			}
			unicornhat.Show()
		}
	})
}

func main() {
	defer unicornsys.Terminate(0)
	flag.Parse()
	neuron.StopRunningOnSignal(syscall.SIGINT)
	if *brightness > unicornhat.DefaultBrightness() {
		if *brightness > 1.0 {
			fmt.Println("Brightness higher than 1.0, setting to default brightness for safety sake!")
			*brightness = unicornhat.DefaultBrightness()
		} else {
			fmt.Println("WARNING: you've set the brightness higher than the default, this can get bright! Please don't look directly at it!")
		}
	}
	if *brightness < 0.0 {
		fmt.Println("Brightness less than 0.0, setting to default brightness!")
		*brightness = unicornhat.DefaultBrightness()
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
