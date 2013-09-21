package main

import (
	"fmt"
	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
	"io"
	"os"
	"path"
	"strings"
)

type AbsoluteResolver struct {
	prefix string
}

func (r *AbsoluteResolver) Resolve(modPath string) (io.Reader, error) {
	if r.prefix == "" {
		fmt.Println("No path set. Using default")
		r.prefix = os.Getenv("GOPATH") + "/src/github.com/PuerkitoBio/agora/runtime/stdlib/"
	}

	toOpen := path.Join(r.prefix, modPath)

	if strings.HasPrefix(modPath, "./") {
		curDir, _ := os.Getwd()
		toOpen = path.Join(curDir, modPath)
	}

	f, err := os.Open(toOpen)

	if err != nil {
		return nil, err
	}

	return f, nil

}

func main() {
	fmt.Println("hey there, about to run some agora")
	ctx := runtime.NewCtx(new(runtime.FileResolver), new(compiler.Compiler))
	f, _ := os.Open("./ago/test.ago")
	defer f.Close()
	ctx.Compiler.Compile("./ago/test.ago", f)
	ctx.RegisterNativeModule(new(stdlib.FmtMod))
	ctx.RegisterNativeModule(new(stdlib.FilepathMod))
	ctx.RegisterNativeModule(new(stdlib.ConvMod))
	ctx.RegisterNativeModule(new(stdlib.StringsMod))
	ctx.RegisterNativeModule(new(stdlib.MathMod))
	ctx.RegisterNativeModule(new(stdlib.OsMod))
	ctx.RegisterNativeModule(new(stdlib.TimeMod))
	mod, err := ctx.Load("./ago/test.ago")
	if err != nil {
		fmt.Println("Couldn't load module", err.Error())
		os.Exit(1)
	}
	val, err := mod.Run()
	if err != nil {
		fmt.Println("Error executing module", err.Error())
		os.Exit(1)
	}
	fmt.Println(val)
	fmt.Println("Done here.")
}
