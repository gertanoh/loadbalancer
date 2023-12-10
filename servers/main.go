package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gertanoh/loadbalancer/internal/server"
)

var port int

func main() {

	flag.IntVar(&port, "port", 8080, "API server port")
	flag.Parse()

	handler := handleRequests()
	var config server.Config
	config.HttpAddr = fmt.Sprintf(":%d", port)
	config.Handler = handler

	srv, err := server.New(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Serve(); err != nil {
		log.Fatal(err)
	}
}

func handleRequests() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("---------------------------")
		log.Println("Received request from ", r.RemoteAddr)
		log.Println(r.Method, r.URL.Path, r.Proto)

		for name, headers := range r.Header {
			for _, h := range headers {
				log.Println(name, h)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		resp := "Replied with hello from server : " + strconv.Itoa(port) + "\n"
		w.Write([]byte(resp))
	})
}
