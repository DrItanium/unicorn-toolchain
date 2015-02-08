package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/DrItanium/unicornhat"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var coords = flag.Bool("horizontal", false, "display lights in a horzontal fashion instead of vertical")
var randomize = flag.Bool("randomize", false, "randomize the position to select")
var hyperSpeed = flag.Bool("hyperspeed", false, "eliminate all microsecond delay calls")
var pixelDelay = flag.Uint("delay", 10, "number of microseconds to pause in between pixel updates")
var preset = flag.Int("preset", 0, "select a predefined color palette index.\n\t0 - all colors\n\t1 - greyscale\n\t2 - green\n\t3 - yellow\n\t4 - yellow and green\n\t5 - purple\n\t6 - cyan")
var fullPage = flag.Bool("fullpage", false, "Update all 64 pixels at a time instead of updating after each pixel change")

const (
	FullColorPreset = iota
	Greyscale
	Green
	Yellow
	GreenAndYellow
	Purple
	Cyan
)

type intensityHandler func(intensity byte) unicornhat.Pixel

func singleBytePixel(input *bufio.Reader, i int, transformer intensityHandler) bool {
	tmp, err := input.ReadByte()
	if err != nil {
		return false
	}
	unicornhat.SetPixelColorType(uint(i), transformer(tmp))
	if !*fullPage {
		unicornhat.Show()
	}
	if !*hyperSpeed {
		microsecond_delay(time.Duration(*pixelDelay))
	}
	return true
}
func terminate_unicorn(status int) {
	for i := 0; i < 64; i++ {
		unicornhat.SetPixelColor(uint(i), 0, 0, 0)
	}
	unicornhat.Show()
	unicornhat.Shutdown(status)
}
func microsecond_delay(usec time.Duration) {
	time.Sleep(usec * time.Microsecond)
}
func colorPixel(input *bufio.Reader, i int) bool {
	var pixel unicornhat.Pixel
	tmp, err := input.ReadByte()
	if err != nil {
		return false
	} else {
		pixel.R = tmp
	}
	tmp, err = input.ReadByte()
	if err != nil {
		return false
	} else {
		pixel.G = tmp
	}
	tmp, err = input.ReadByte()
	if err != nil {
		return false
	} else {
		pixel.B = tmp
	}
	unicornhat.SetPixelColorType(uint(i), pixel)
	if !*fullPage {
		unicornhat.Show()
	}
	if !*hyperSpeed {
		microsecond_delay(time.Duration(*pixelDelay))
	}
	return true
}
func greyscalePixel(input *bufio.Reader, i int) bool {
	return singleBytePixel(input, i, func(intensity byte) unicornhat.Pixel {
		return unicornhat.Pixel{R: intensity, G: intensity, B: intensity}
	})
}
func purplePixel(input *bufio.Reader, i int) bool {
	// purple is made up of red and blue so no green
	return singleBytePixel(input, i, func(intensity byte) unicornhat.Pixel { return unicornhat.Pixel{R: intensity, G: 0, B: intensity} })
}
func cyanPixel(input *bufio.Reader, i int) bool {
	// cyan is made up of green and blue so no red
	return singleBytePixel(input, i, func(intensity byte) unicornhat.Pixel { return unicornhat.Pixel{R: 0, G: intensity, B: intensity} })
}
func yellowPixel(input *bufio.Reader, i int) bool {
	// yellow is made up of red and green so no blue
	return singleBytePixel(input, i, func(intensity byte) unicornhat.Pixel { return unicornhat.Pixel{R: intensity, G: intensity, B: 0} })
}
func greenPixel(input *bufio.Reader, i int) bool {
	return singleBytePixel(input, i, func(intensity byte) unicornhat.Pixel { return unicornhat.Pixel{R: 0, G: intensity, B: 0} })
}
func showPixel(input *bufio.Reader, i int) bool {
	switch *preset {
	case FullColorPreset:
		return colorPixel(input, i)
	case Greyscale:
		return greyscalePixel(input, i)
	case Green:
		return greenPixel(input, i)
	case Yellow:
		return yellowPixel(input, i)
	case GreenAndYellow:
		if rand.Int()%2 == 1 {
			return greenPixel(input, i)
		} else {
			return yellowPixel(input, i)
		}
	case Purple:
		return purplePixel(input, i)
	case Cyan:
		return cyanPixel(input, i)
	default:
		fmt.Println("Invalid preset ", *preset, "provided!")
		return false
	}
}

func main() {
	defer terminate_unicorn(0)
	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	running := true
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		running = false
	}()
	unicornhat.Initialize(64, unicornhat.DefaultBrightness())
	unicornhat.ClearLEDBuffer()
	unicornhat.Show()
	input := bufio.NewReader(os.Stdin)
	// setup the initial pixels
	for running {
		if *coords {
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					vY := y
					vX := x
					if *randomize && (r.Int()%2 == 1) {
						vY = r.Int() % 8
						vX = r.Int() % 8
					}
					running = showPixel(input, unicornhat.CoordinateToPosition(vX, vY)) && running
					if !running {
						break
					}
				}
				if !running {
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
				running = showPixel(input, vI) && running
				if !running {
					break
				}
			}
			if !running {
				break
			}
			if *fullPage {
				unicornhat.Show()
			}
		}
	}
}
