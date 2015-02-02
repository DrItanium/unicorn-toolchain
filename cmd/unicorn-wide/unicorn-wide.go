package main

import "github.com/DrItanium/unicornhat"
import "time"
import "math/rand"
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
func random_byte() byte {
	return byte(rand.Int())
}

func main() {
	running := true
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		running = false
	}()
	unicornhat.Initialize(64, unicornhat.DefaultBrightness())
	unicornhat.ClearLEDBuffer()
	for running {
		unicornhat.SetPixelColor(0, random_byte(), random_byte(), random_byte())
		unicornhat.Show()
		microsecond_delay(10)
		for index := 1; index < 64; index++ {
			unicornhat.SetPixelColor(uint(index), random_byte(), random_byte(), random_byte())
			for i := 0; i < index; i++ {
				pix := unicornhat.GetPixelColor(uint(i))
				var r, g, b byte
				if pix.R > 0 {
					r = pix.R - 1
				} else {
					r = 0
				}
				if pix.G > 0 {
					g = pix.G - 1
				} else {
					g = 0
				}
				if pix.B > 0 {
					b = pix.B - 1
				} else {
					b = 0
				}
				unicornhat.SetPixelColor(uint(i), r, g, b)
			}
			unicornhat.Show()
			microsecond_delay(10)
		}
		moreToFade := true
		for moreToFade {

			moreToFade = false
			for i := 0; i < 64; i++ {
				pix := unicornhat.GetPixelColor(uint(i))
				var r, g, b byte
				if pix.R > 0 {
					r = pix.R - 1
					moreToFade = true
				} else {
					r = 0
				}
				if pix.G > 0 {
					g = pix.G - 1
					moreToFade = true
				} else {
					g = 0
				}
				if pix.B > 0 {
					b = pix.B - 1
					moreToFade = true
				} else {
					b = 0
				}
				unicornhat.SetPixelColor(uint(i), r, g, b)
			}
			unicornhat.Show()
		}
		microsecond_delay(10)
	}
	defer terminate_unicorn(0)
}
