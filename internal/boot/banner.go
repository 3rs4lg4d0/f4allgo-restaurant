package boot

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mbndr/figlet4go"
)

func PrintBanner() {
	if GetConfig().AppBanner {
		fmt.Printf("Booting '%s' at %s\n\n", GetConfig().AppName, time.Now().Format("2006-01-02T15:04:05-07:00"))
		asciiFont := figlet4go.NewAsciiRender()
		renderStr, _ := asciiFont.Render(GetConfig().AppName)
		fmt.Print(renderStr)
		fmt.Printf("\nGo version: %s\n\n", runtime.Version())
	}
}
