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

var fullpageUpdate = flag.Bool("fullpage", false, "Perform full page updates")
var columnByColumn = flag.Bool("column", false, "Perform column by column updates")
var rowByRow = flag.Bool("row", false, "Perform row by row updates")
var hyperspeed = flag.Bool("hyperspeed", false, "Disable delay (still consumes delay bytes)")
var brightness = flag.Float64("brightness-factor", unicornhat.DefaultBrightness(), "Set brightness cap (0.0 - 1.0).\n\tWARNING: If you set this brightness too high you can cause retinal damage and I'm not responsible for that!!!")

var running = true

func terminate_unicorn(status int) {
	for i := 0; i < 64; i++ {
		unicornhat.SetPixelColor(uint(i), 0, 0, 0)
	}
	unicornhat.Show()
	unicornhat.Shutdown(status)
}
func microsecond_delay(usec time.Duration) {
	if !*hyperspeed {
		time.Sleep(usec * time.Microsecond)
	}
}
func random_byte() byte {
	return byte(rand.Int())
}

type MicrocodeField struct {
	Pixel unicornhat.Pixel
	Delay byte
}
type MicrocodeWord [64]MicrocodeField

func readMicrocodeWord(input *bufio.Reader) (*MicrocodeWord, error) {
	elements := make([]byte, 256)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		running = false
		return nil, err
	}
	var word MicrocodeWord
	for i := 0; i < 256; i += 4 {
		word[i/4].Pixel.R = elements[i]
		word[i/4].Pixel.G = elements[i+1]
		word[i/4].Pixel.B = elements[i+2]
		word[i/4].Delay = elements[i+3]
	}
	return &word, nil
}

func pixelByPixelUpdate(input *bufio.Reader) {
	word, err := readMicrocodeWord(input)
	if err != nil {
		return
	}
	for i := 0; i < 64; i++ {
		if !running {
			return
		}
		unicornhat.SetPixelColorType(uint(i), &word[i].Pixel)
		unicornhat.Show()
		microsecond_delay(time.Duration(word[i].Delay))
	}
}

func fullPageUpdate(input *bufio.Reader) {
	word, err := readMicrocodeWord(input)
	if err != nil {
		return
	}
	for i := 0; i < 64; i++ {
		if !running {
			return
		}
		unicornhat.SetPixelColorType(uint(i), &word[i].Pixel)
		microsecond_delay(time.Duration(word[i].Delay))
	}
	unicornhat.Show()
}

func columnByColumnUpdate(input *bufio.Reader) {
	word, err := readMicrocodeWord(input)
	if err != nil {
		return
	}
	for x := 0; x < 8; x++ {
		if !running {
			return
		}
		for y := 0; y < 8; y++ {
			if !running {
				return
			}
			index := unicornhat.CoordinateToPosition(x, y)
			unicornhat.SetPixelColorType(uint(index), &word[index].Pixel)
			microsecond_delay(time.Duration(word[index].Delay))
		}
		unicornhat.Show()
	}
}

func rowByRowUpdate(input *bufio.Reader) {
	word, err := readMicrocodeWord(input)
	if err != nil {
		return
	}
	for y := 0; y < 8; y++ {
		if !running {
			return
		}
		for x := 0; x < 8; x++ {
			if !running {
				return
			}
			index := unicornhat.CoordinateToPosition(x, y)
			unicornhat.SetPixelColorType(uint(index), &word[index].Pixel)
			microsecond_delay(time.Duration(word[index].Delay))
		}
		unicornhat.Show()
	}
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
	if *brightness > unicornhat.DefaultBrightness() {
		if *brightness > 1.0 {
			fmt.Println("Brightness higher than 1.0, setting to default brightness for safety sake!")
			*brightness = unicornhat.DefaultBrightness()
		} else {
			fmt.Println("WARNING: you've set the brightness higher than the default, this can get bright! Please don't look directly at it!")
		}
	}
	if *brightness < 0.0 {
		fmt.Println("Brightness less than 0.0, setting to default brightness!")
		*brightness = unicornhat.DefaultBrightness()
	}
	unicornhat.Initialize(64, *brightness)
	unicornhat.ClearLEDBuffer()
	input := bufio.NewReader(os.Stdin)
	for running {
		if *fullpageUpdate {
			fullPageUpdate(input)
		} else if *columnByColumn {
			columnByColumnUpdate(input)
		} else if *rowByRow {
			rowByRowUpdate(input)
		} else {
			pixelByPixelUpdate(input)
		}
	}
}
