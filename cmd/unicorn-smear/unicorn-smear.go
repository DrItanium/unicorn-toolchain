package main

import "github.com/DrItanium/unicornhat"
import "time"

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
	unicornhat.SetBrightness(unicornhat.DefaultBrightness())
	unicornhat.Init(64)
	unicornhat.ClearLEDBuffer()
	for count := 0; count < 64; count++ {
		unicornhat.SetPixelColor(0, 255, 255, 255)
		unicornhat.Show()
		microsecond_delay(10)
		for index := 1; index < 64; index++ {
			unicornhat.SetPixelColor(uint(index), 255, 255, 255)
			for prev := 0; prev < index; prev++ {
				pix := unicornhat.GetPixelColor(uint(prev))
				unicornhat.SetPixelColor(uint(prev), pix.R-1, pix.G-1, pix.B-1)
			}
			unicornhat.Show()
			microsecond_delay(10)
		}
		shouldContinue := true
		for shouldContinue {
			shouldContinue = false
			for i := 0; i < 64; i++ {
				pix := unicornhat.GetPixelColor(uint(i))
				var r, g, b byte
				if pix.R > 0 {
					r = pix.R - 1
					shouldContinue = true
				} else {
					r = 0
				}
				if pix.G > 0 {
					g = pix.G - 1
					shouldContinue = true
				} else {
					g = 0
				}
				if pix.B > 0 {
					b = pix.B - 1
					shouldContinue = true
				} else {
					b = 0
				}
				unicornhat.SetPixelColor(uint(i), r, g, b)
			}
			unicornhat.Show()
			microsecond_delay(10)
		}
		microsecond_delay(10)
	}
	defer terminate_unicorn(0)
}
