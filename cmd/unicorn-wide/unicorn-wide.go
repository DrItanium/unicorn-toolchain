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
func updateColumn(input *bufio.Reader, column int) {
	elements := make([]byte, 32)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		running = false
		return
	}
	var pixel unicornhat.Pixel
	for i := 0; i < 32; i += 4 {
		pixel.R = elements[i]
		pixel.G = elements[i+1]
		pixel.B = elements[i+2]
		delay := elements[i+3]
		unicornhat.SetPixelColorType(uint(unicornhat.CoordinateToPosition(column, int(i/4))), &pixel)
		microsecond_delay(time.Duration(delay))
	}
	unicornhat.Show()
}

func updateRow(input *bufio.Reader, row int) {
	elements := make([]byte, 32)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		running = false
		return
	}
	var pixel unicornhat.Pixel
	for i := 0; i < 32; i += 4 {
		pixel.R = elements[i]
		pixel.G = elements[i+1]
		pixel.B = elements[i+2]
		delay := elements[i+3]
		unicornhat.SetPixelColorType(uint(unicornhat.CoordinateToPosition(int(i/4), row)), &pixel)
		microsecond_delay(time.Duration(delay))
	}
	unicornhat.Show()
}

func columnByColumnUpdate(input *bufio.Reader) {
	for x := 0; x < 8; x++ {
		updateColumn(input, x)
		if !running {
			return
		}
	}
}

func rowByRowUpdate(input *bufio.Reader) {
	for y := 0; y < 8; y++ {
		updateRow(input, y)
		if !running {
			return
		}
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
