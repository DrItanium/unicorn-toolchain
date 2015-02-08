package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/DrItanium/neuron"
	"github.com/DrItanium/unicorn-toolchain/sys"
	"github.com/DrItanium/unicornhat"
	"syscall"
	"time"
)

var horizontal = flag.Bool("horizontal", false, "display lights in a horzontal fashion instead of vertical")
var randomize = flag.Bool("randomize", false, "randomize the position to select")
var hyperSpeed = flag.Bool("hyperspeed", false, "eliminate all microsecond delay calls")
var pixelDelay = flag.Uint("delay", 10, "number of microseconds to pause in between pixel updates")
var preset = flag.Uint("preset", 0, "select a predefined color palette index.\n\t0 - all colors\n\t1 - greyscale\n\t2 - red\n\t3 - green\n\t4 - blue\n\t5 - yellow\n\t6 - purple\n\t7 - cyan")
var combineWith = flag.Uint("combine-with", 0, "select another preset to alternate with (0 disables the feature)")
var fullPage = flag.Bool("fullpage", false, "Update all 64 pixels at a time instead of updating after each pixel change")

const (
	FullColorPreset = iota
	Greyscale
	Red
	Green
	Blue
	Yellow
	Purple
	Cyan
)

type intensityHandler func(intensity byte) *unicornhat.Pixel
type colorTransfomer func(intensities []byte) *unicornhat.Pixel

var intensityFuncs = []intensityHandler{
	FullColorPreset: nil,
	Greyscale:       func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(intensity, intensity, intensity) },
	Red:             func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(intensity, 0, 0) },
	Green:           func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(0, intensity, 0) },
	Blue:            func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(0, 0, intensity) },
	Yellow:          func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(intensity, intensity, 0) },
	Purple:          func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(intensity, 0, intensity) },
	Cyan:            func(intensity byte) *unicornhat.Pixel { return unicornhat.NewPixel(0, intensity, intensity) },
}

func dualIntensity(primary, secondary uint) intensityHandler {
	return func(intensity byte) *unicornhat.Pixel {
		if neuron.GlobalRandomBool() {
			return intensityFuncs[primary](intensity)
		} else {
			return intensityFuncs[secondary](intensity)
		}
	}
}

func hasValidDualIntensity(primary, secondary uint) bool {
	return (primary != secondary) && primary != 0 && secondary != 0
}

func microsecond_delay(usec time.Duration) {
	if !*hyperSpeed {
		time.Sleep(usec * time.Microsecond)
	}
}
func updatePixel(input *bufio.Reader, i int, sz int, transform colorTransfomer) {
	elements := make([]byte, sz)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		neuron.StopRunning()
		return
	}
	unicornhat.SetPixelColorType(uint(i), transform(elements))
	microsecond_delay(time.Duration(*pixelDelay))
	if !*fullPage {
		unicornhat.Show()
	}
}
func singleBytePixel(input *bufio.Reader, i int, transformer intensityHandler) {
	updatePixel(input, i, 1, func(elements []byte) *unicornhat.Pixel {
		return transformer(elements[0])
	})
}
func colorPixel(input *bufio.Reader, i int) {
	updatePixel(input, i, 3, func(elements []byte) *unicornhat.Pixel {
		return unicornhat.NewPixel(elements[0], elements[1], elements[2])
	})
}

func showPixel(input *bufio.Reader, i int) {
	if *preset >= uint(len(intensityFuncs)) {
		neuron.StopRunning()
		fmt.Println("Invalid preset", *preset, "provided!")
	} else if *combineWith >= uint(len(intensityFuncs)) {
		neuron.StopRunning()
		fmt.Println("Invalid combineWith value", *combineWith, "provided!")
	} else {
		if hasValidDualIntensity(*preset, *combineWith) {
			singleBytePixel(input, i, dualIntensity(*preset, *combineWith))
		} else if *preset == FullColorPreset {
			colorPixel(input, i)
		} else {
			singleBytePixel(input, i, intensityFuncs[*preset])
		}
	}
}

func main() {
	defer unicornsys.Terminate(0)
	flag.Parse()
	r := neuron.NewTimeSourcedRand()
	input := neuron.NewStandardInReader()
	neuron.StopRunningOnSignal(syscall.SIGINT)
	unicornhat.Initialize(64, unicornhat.DefaultBrightness())
	unicornhat.ClearLEDBuffer()
	for neuron.IsRunning() {
		for y := 0; neuron.IsRunning() && y < 8; y++ {
			for x := 0; neuron.IsRunning() && x < 8; x++ {
				vY := y
				vX := x
				if *randomize && neuron.RandomBool(r) {
					vY = r.Int() % 8
					vX = r.Int() % 8
				}
				var index int
				if *horizontal {
					index = unicornhat.CoordinateToPosition(vX, vY)
				} else {
					index = unicornhat.CoordinateToPosition(vY, vX)
				}
				showPixel(input, index)
			}
		}
		if *fullPage {
			unicornhat.Show()
		}
	}
}
