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

var fullpageUpdate = flag.Bool("fullpage", false, "Perform full page updates instead of pixel by pixel")

//var rowbyrowUpdate = flag.Bool("rowbyrow", false, "Peform
var running = true

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
func random_byte() byte {
	return byte(rand.Int())
}

func readPixelByPixelQuantum(input *bufio.Reader) (*unicornhat.Pixel, error) {
	elements := make([]byte, 4)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		return nil, err
	}
	microsecond_delay(time.Duration(elements[3]))
	return unicornhat.NewPixel(elements[0], elements[1], elements[2]), err
}

func pixelByPixelUpdate(input *bufio.Reader) {
	for i := 0; i < 64; i++ {
		if !running {
			return
		}
		pixel, err := readPixelByPixelQuantum(input)
		if err != nil {
			running = false
		}
		if pixel != nil {
			unicornhat.SetPixelColorType(uint(i), pixel)
		}
		unicornhat.Show()
	}
}

func fullPageUpdate(input *bufio.Reader) {
	elements := make([]byte, 256)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		running = false
		return
	}
	var pixel unicornhat.Pixel
	for i := 0; i < 256; i += 4 {
		pixel.R = elements[i]
		pixel.G = elements[i+1]
		pixel.B = elements[i+2]
		delay := elements[i+3]
		unicornhat.SetPixelColorType(uint(i/4), &pixel)
		microsecond_delay(time.Duration(delay))
	}
	unicornhat.Show()
}

func main() {
	defer terminate_unicorn(0)
	flag.Parse()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		running = false
	}()
	unicornhat.Initialize(64, unicornhat.DefaultBrightness())
	unicornhat.ClearLEDBuffer()
	input := bufio.NewReader(os.Stdin)
	for running {
		if *fullpageUpdate {
			fullPageUpdate(input)
		} else {
			pixelByPixelUpdate(input)
		}
	}
}
