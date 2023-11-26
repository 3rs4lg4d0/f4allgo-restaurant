package boot

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mbndr/figlet4go"
)

func PrintBanner(additionalInfo ...string) {
	if GetConfig().AppBanner {
		fmt.Printf("Booting '%s' at %s\n\n", GetConfig().AppName, time.Now().Format("2006-01-02T15:04:05-07:00"))

		asciiFont := figlet4go.NewAsciiRender()
		renderStr, _ := asciiFont.Render(GetConfig().AppName)
		fmt.Print(renderStr)
		fmt.Printf("\nGo version: %s\n", runtime.Version())
		for _, ai := range additionalInfo {
			fmt.Printf("%s\n", ai)
		}
	}
}
