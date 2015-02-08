package main

import (
	"bufio"
	"github.com/DrItanium/unicornhat"
	"math/rand"
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
func random_byte() byte {
	return byte(rand.Int())
}

func readPixelByPixelQuantum(input *bufio.Reader) (*unicornhat.Pixel, error) {
	elements := make([]byte, 3)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		return nil, err
	}
	return &unicornhat.Pixel{R: elements[0], G: elements[1], B: elements[2]}, err
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
	input := bufio.NewReader(os.Stdin)
	for running {
		for i := 0; i < 64; i++ {
			if !running {
				break
			}
			pixel, err := readPixelByPixelQuantum(input)
			if err != nil {
				running = false
			}
			if pixel != nil {
				unicornhat.SetPixelColorType(uint(i), pixel)
			}
		}
	}
}
