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
	// setup the initial pixels
	for running {
		var pixel unicornhat.Pixel
		for r := 0; r < 255; r += 5 {
			pixel.R = byte(r)
			for g := 0; g < 255; g += 5 {
				pixel.G = byte(g)
				for b := 0; b < 255; b += 5 {
					pixel.B = byte(b)
					for i := 0; i < 64; i++ {
						unicornhat.SetPixelColorType(uint(i), pixel)
						if !running {
							break
						}
					}
					if !running {
						break
					}
					unicornhat.Show()
				}
				if !running {
					break
				}
			}
			if !running {
				break
			}
		}
		if !running {
			break
		}
	}
}
