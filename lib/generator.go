package libLiveSrv

import (
	"fmt"
	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
	"os"
)

// NewAgoraClosure loads a full agora context with the module passed in, then
// returns a closure that calls the 'Run' method on the closure and returns the
// result
func NewAgoraClosure(modPath string) func() string {
	ctx := runtime.NewCtx(new(HTTPResolver), new(compiler.Compiler))
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
