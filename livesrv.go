package main

import (
	"fmt"
	"github.com/AeroNotix/wedge"
	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
	"net/http"
	"os"
	"time"
)

func NewAgoraClosure(modPath string) func() string {
	ctx := runtime.NewCtx(new(runtime.FileResolver), new(compiler.Compiler))
	f, _ := os.Open(modPath)
	defer f.Close()

	ctx.Compiler.Compile(modPath, f)

	ctx.RegisterNativeModule(new(stdlib.FmtMod))
	ctx.RegisterNativeModule(new(stdlib.FilepathMod))
	ctx.RegisterNativeModule(new(stdlib.ConvMod))
	ctx.RegisterNativeModule(new(stdlib.StringsMod))
	ctx.RegisterNativeModule(new(stdlib.MathMod))
	ctx.RegisterNativeModule(new(stdlib.OsMod))
	ctx.RegisterNativeModule(new(stdlib.TimeMod))

	mod, err := ctx.Load(modPath)
	if err != nil {
		fmt.Println("Couldn't load module", err.Error())
		os.Exit(1)
	}

	return func() string {

		val, err := mod.Run()
		if err != nil {
			fmt.Println("Error executing module", err.Error())
			os.Exit(1)
		}

		return val.String()
	}
}

var ago = NewAgoraClosure("./ago/test.ago")

func index(w http.ResponseWriter, r *http.Request) (string, int) {
	return ago(), 200
}

func main() {
	app := wedge.NewAppServer("8080", time.Second*2)

	app.AddURLs(
		wedge.URL("^/$", "index", index, 1),
	)

	go func() {
		<-time.After(time.Second * 5)
		fmt.Println("Swapping context")
		ago = NewAgoraClosure("./ago/test2.ago")
	}()

	app.Run()
}
