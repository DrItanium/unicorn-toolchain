// microcode structures
package microcode

import (
	"bufio"
	"github.com/DrItanium/neuron"
	"github.com/DrItanium/unicornhat"
	"time"
)

const (
	MicrocodeWordSize   = 256
	MicrocodeFieldCount = 64
)

type MicrocodeField struct {
	Pixel unicornhat.Pixel
	Delay byte
}
type MicrocodeWord [MicrocodeFieldCount]MicrocodeField

func ReadMicrocodeWord(input *bufio.Reader) (*MicrocodeWord, error) {
	elements := make([]byte, MicrocodeWordSize)
	count, err := input.Read(elements)
	if err != nil && count == 0 {
		neuron.StopRunning()
		return nil, err
	}
	var word MicrocodeWord
	for i := 0; i < MicrocodeWordSize; i += 4 {
		index := i / 4
		word[index].Pixel.R = elements[i]
		word[index].Pixel.G = elements[i+1]
		word[index].Pixel.B = elements[i+2]
		word[index].Delay = elements[i+3]
	}
	return &word, nil
}

type DelayFunction func(delay time.Duration)

func (self *MicrocodeField) Pause(f DelayFunction) {
	f(time.Duration(self.Delay))
}

func (self *MicrocodeField) UpdateNativePixel(index int, f DelayFunction) {
	unicornhat.SetPixelColorType(uint(index), &self.Pixel)
	self.Pause(f)
}

func (self *MicrocodeWord) UpdateNativePixels(f DelayFunction) {
	for i := 0; i < MicrocodeFieldCount; i++ {
		self[i].UpdateNativePixel(i, f)
	}
}

func (self *MicrocodeWord) FieldFromCoordinates(x, y int) *MicrocodeField {
	return &self[unicornhat.CoordinateToPosition(x, y)]
}
