// microcode structures
package microcode

import (
	"bufio"
	"github.com/DrItanium/neuron"
	"github.com/DrItanium/unicornhat"
)

const (
	MicrocodeWordSize = 256
)

type MicrocodeField struct {
	Pixel unicornhat.Pixel
	Delay byte
}
type MicrocodeWord [64]MicrocodeField

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
