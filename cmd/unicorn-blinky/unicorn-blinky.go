package main

import "github.com/DrItanium/unicornhat"
import "time"
import "fmt"

func terminate_unicorn(status int) {
	for i := 0; i < 64; i++ {
		unicornhat.SetPixelColor(uint(i),0,0,0)
	}
	unicornhat.Show()
	unicornhat.Shutdown(status)
}
func main() {
	unicornhat.SetBrightness(unicornhat.DefaultBrightness())
	unicornhat.Init(64)
	unicornhat.ClearLEDBuffer()
	for count := 0; count < 64; count++ {
		fmt.Println("count = ", count)
		unicornhat.SetPixelColor(0,1,1,1)
		unicornhat.Show()
		time.Sleep(10 * time.Microsecond)
		for index := 1; index < 64; index++ {
			unicornhat.SetPixelColor(uint(index - 1), 0, 0, 0)
			unicornhat.SetPixelColor(uint(index), 255, 255, 255)
			unicornhat.Show()
			time.Sleep(10)
		}
		unicornhat.SetPixelColor(63,0,0,0)
		unicornhat.Show()
		time.Sleep(10 * time.Microsecond)
	}
	defer terminate_unicorn(0)
}
