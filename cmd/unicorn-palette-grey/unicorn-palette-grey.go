package main

import (
	"bufio"
	"github.com/DrItanium/unicornhat"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
func main() {
	defer terminate_unicorn(0)
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
		for i := 0; i < 64; i++ {
			var pixel unicornhat.Pixel
			tmp, err := input.ReadByte()
			if err != nil {
				running = false
			} else {
				pixel.R = tmp
			}
			tmp, err = input.ReadByte()
			if err != nil {
				running = false
			} else {
				pixel.G = tmp
			}
			tmp, err = input.ReadByte()
			if err != nil {
				running = false
			} else {
				pixel.B = tmp
			}
			unicornhat.SetPixelColorType(uint(i), pixel)
			unicornhat.Show()
			microsecond_delay(10)
			if !running {
				break
			}
		}
	}
}
