package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amukherj/envoys2s/internal/rand"
	"github.com/amukherj/envoys2s/internal/server"
)

const (
	primeDiv = 823553
)

type argSet struct {
	lstfile *string
	port    *uint
}

func handleArgs() (*argSet, error) {
	lstfile := flag.String("list", "", "Path to file containing item data")
	port := flag.Uint("port", 0, "Port to listen to")
	flag.Parse()

	if *lstfile == "" {
		return nil, fmt.Errorf("No item file specified")
	}

	if *port == 0 {
		return nil, fmt.Errorf("Port cannot be zero")
	}

	return &argSet{
		lstfile: lstfile,
		port:    port,
	}, nil
}

func readItems(lstfile string) ([]string, error) {
	file, err := os.Open(lstfile)
	if err != nil {
		return nil, fmt.Errorf("Could not open file %s: %v", lstfile, err)
	}

	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Could not read file %s: %v\n", lstfile, err)
	}

	str := string(buf)
	items := strings.Split(str, "\n")

	for n, item := range items {
		if len(item) == 0 {
			fmt.Fprintf(os.Stderr, "Empty at %d\n", n)
		}
	}

	return items, nil
}

func main() {
	args, err := handleArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	items, err := readItems(*args.lstfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	rg := rand.NewRandGen(time.Now().UnixNano()%primeDiv, len(items)-1)

	routes := map[string]http.HandlerFunc{}
	routes["/"] = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			item := items[rg.NextInt()]
			log.Printf("Served request with %s.", item)
			io.WriteString(w, item)
		})

	routes["/es"] = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Served ES request.")
			io.WriteString(w, "hola, mundo!\n")
		})

	svr := server.NewServer(*args.port)
	chnl := svr.Start(routes)
	<-chnl
	svr.Stop()
}
