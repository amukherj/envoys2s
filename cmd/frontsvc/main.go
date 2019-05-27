package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"envoys2s/internal/httpipe"
	"envoys2s/internal/server"
)

var (
	items = []string{}
)

func main() {
	port := flag.Uint("port", 0, "Port to listen to")
	flag.Parse()

	if *port == 0 {
		fmt.Fprintf(os.Stderr, "Port can't be zero\n")
		return
	}

	ctx, _ := context.WithCancel(context.Background())
	time.Sleep(2 * time.Second)

	// For header-based routing, remove the moods suffix if any in the URL
	// For prefix-bassed routing, remove the header(s)
	src1 := httpipe.NewPipe(ctx, "http://localhost:9001/moods",
		map[string]string{
			/* "x-target": "moods", */
		})

	// For header-based routing, remove the rockers suffix if any in the URL
	// For prefix-bassed routing, remove the header(s)
	src2 := httpipe.NewPipe(ctx, "http://localhost:9001/rockers",
		map[string]string{
			/* "x-target": "rockers", */
		})

	routes := map[string]http.HandlerFunc{}
	routes["/"] = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r1 := <-src1
			r2 := <-src2
			if r1.Err == nil && r2.Err == nil {
				io.WriteString(w, string(r1.Result)+" "+string(r2.Result))
			} else {
				io.WriteString(w, fmt.Sprintf("%s(%v) %s(%v)",
					r1.Result, r1.Err, r2.Result, r2.Err))
			}
		})

	routes["/es"] = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "hola, mundo!\n")
		})

	svr := server.NewServer(*port)
	chnl := svr.Start(routes)
	defer svr.Stop()
	<-chnl
}
