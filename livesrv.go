package main

import (
	"fmt"
	"github.com/AeroNotix/wedge"
	"github.com/ryansb/livesrv/lib"
	"net/http"
	"time"
)

var ago = libLiveSrv.NewAgoraClosure("./ago/test.ago")

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
		ago = libLiveSrv.NewAgoraClosure("./ago/test2.ago")
	}()

	app.Run()
}
