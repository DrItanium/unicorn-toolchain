// general unicorn hat system functions
package unicornsys

import (
	"github.com/DrItanium/unicornhat"
)

func Terminate(status int) {
	for i := 0; i < 64; i++ {
		unicornhat.SetPixelColor(uint(i), 0, 0, 0)
	}
	unicornhat.Show()
	unicornhat.Shutdown(status)
}
