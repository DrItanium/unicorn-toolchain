package main

import (
	"bufio"
	"flag"
	"github.com/DrItanium/unicornhat"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var coords = flag.Bool("horizontal", false, "display lights in a horzontal fashion instead of vertical")
var randomize = flag.Bool("randomize", false, "randomize the position to select")

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
func showPixel(input *bufio.Reader, i int) bool {
	var pixel unicornhat.Pixel
	tmp, err := input.ReadByte()
	if err != nil {
		return false
	} else {
		pixel.R = tmp
		pixel.G = tmp
		pixel.B = tmp
	}
	unicornhat.SetPixelColorType(uint(i), pixel)
	unicornhat.Show()
	microsecond_delay(100)
	return true
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
		}
	}
}
