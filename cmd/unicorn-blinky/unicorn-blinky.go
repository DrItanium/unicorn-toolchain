package main

import "github.com/DrItanium/unicornhat"
import "time"
import "os"
import "os/signal"
import "syscall"

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
	running := true
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		running = false
	}()
	unicornhat.SetBrightness(unicornhat.DefaultBrightness())
	unicornhat.Init(64)
	unicornhat.ClearLEDBuffer()
	for running {
		unicornhat.SetPixelColor(0, 1, 1, 1)
		unicornhat.Show()
		time.Sleep(10 * time.Microsecond)
		for index := 1; index < 64; index++ {
			unicornhat.SetPixelColor(uint(index-1), 0, 0, 0)
			unicornhat.SetPixelColor(uint(index), 255, 255, 255)
			unicornhat.Show()
			microsecond_delay(10)
		}
		unicornhat.SetPixelColor(63, 0, 0, 0)
		unicornhat.Show()
		microsecond_delay(10)
	}
	defer terminate_unicorn(0)
}
