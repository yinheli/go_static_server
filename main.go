package main

import (
	"flag"
	"log"
	"net/http"
	"fmt"
)

var (
	// flags

	port = flag.String("p", ":8080", "port")
	dir  = flag.String("d", "", "dir")
)

func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()

	flag.Parse()

	if !flag.Parsed() || *dir == "" {
		flag.Usage()
		return
	}

	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/", maxAgeHandler(0, fs))
	http.ListenAndServe(*port, nil)
}
