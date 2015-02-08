package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/DrItanium/neuron"
	"github.com/DrItanium/unicorn-toolchain/sys"
	"github.com/DrItanium/unicornhat"
	"math/rand"
	"os"
	"syscall"
	"time"
)

var coords = flag.Bool("horizontal", false, "display lights in a horzontal fashion instead of vertical")
var randomize = flag.Bool("randomize", false, "randomize the position to select")
var hyperSpeed = flag.Bool("hyperspeed", false, "eliminate all microsecond delay calls")
var pixelDelay = flag.Uint("delay", 10, "number of microseconds to pause in between pixel updates")
var preset = flag.Int("preset", 0, "select a predefined color palette index.\n\t0 - all colors\n\t1 - greyscale\n\t2 - green\n\t3 - yellow\n\t4 - yellow and green\n\t5 - purple\n\t6 - cyan\n\t7 - blue\n\t8 - red")
var fullPage = flag.Bool("fullpage", false, "Update all 64 pixels at a time instead of updating after each pixel change")

const (
	FullColorPreset = iota
	Greyscale
	Green
	Yellow
	GreenAndYellow
	Purple
	Cyan
	Blue
	Red
)

type intensityHandler func(intensity byte) *unicornhat.Pixel
type colorTransfomer func(intensities []byte) *unicornhat.Pixel

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

func greyscale(intensity byte) *unicornhat.Pixel {
	return unicornhat.NewPixel(intensity, intensity, intensity)
}
func purple(intensity byte) *unicornhat.Pixel {
	// purple is made up of red and blue so no green
	return unicornhat.NewPixel(intensity, 0, intensity)
}
func cyan(intensity byte) *unicornhat.Pixel {
	// cyan is made up of green and blue so no red
	return unicornhat.NewPixel(0, intensity, intensity)
}
func yellow(intensity byte) *unicornhat.Pixel {
	// yellow is made up of green and red so no blue
	return unicornhat.NewPixel(intensity, intensity, 0)
}
func green(intensity byte) *unicornhat.Pixel {
	return unicornhat.NewPixel(0, intensity, 0)
}
func red(intensity byte) *unicornhat.Pixel {
	return unicornhat.NewPixel(intensity, 0, 0)
}

func blue(intensity byte) *unicornhat.Pixel {
	return unicornhat.NewPixel(0, 0, intensity)
}

func showPixel(input *bufio.Reader, i int) {
	switch *preset {
	case FullColorPreset:
		colorPixel(input, i)
	case Greyscale:
		singleBytePixel(input, i, greyscale)
	case Green:
		singleBytePixel(input, i, green)
	case Yellow:
		singleBytePixel(input, i, yellow)
	case GreenAndYellow:
		if rand.Int()%2 == 1 {
			singleBytePixel(input, i, green)
		} else {
			singleBytePixel(input, i, yellow)
		}
	case Purple:
		singleBytePixel(input, i, purple)
	case Cyan:
		singleBytePixel(input, i, cyan)
	case Blue:
		singleBytePixel(input, i, blue)
	case Red:
		singleBytePixel(input, i, red)
	default:
		neuron.StopRunning()
		fmt.Println("Invalid preset ", *preset, "provided!")
	}
}

func main() {
	defer unicornsys.Terminate(0)
	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	neuron.StopRunningOnSignal(syscall.SIGINT)
	unicornhat.Initialize(64, unicornhat.DefaultBrightness())
	unicornhat.ClearLEDBuffer()
	unicornhat.Show()
	input := bufio.NewReader(os.Stdin)
	// setup the initial pixels
	for neuron.IsRunning() {
		if *coords {
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					vY := y
					vX := x
					if *randomize && (r.Int()%2 == 1) {
						vY = r.Int() % 8
						vX = r.Int() % 8
					}
					showPixel(input, unicornhat.CoordinateToPosition(vX, vY))
					if !neuron.IsRunning() {
						break
					}
				}
				if !neuron.IsRunning() {
					break
				}
			}
			if *fullPage {
				unicornhat.Show()
			}
		} else {
			for i := 0; i < 64; i++ {
				vI := i
				if *randomize && (r.Int()%2 == 1) {
					vI = r.Int() % 64
				}
				showPixel(input, vI)
				if !neuron.IsRunning() {
					break
				}
			}
			if !neuron.IsRunning() {
				break
			}
			if *fullPage {
				unicornhat.Show()
			}
		}
	}
}
