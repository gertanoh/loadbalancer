package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

var port int

func main() {

	flag.IntVar(&port, "port", 0, "API server port")
	flag.Parse()

	http.HandleFunc("/", handleRequests)
	log.Println("Load Balancer on port ", port)
	addr := ":" + strconv.Itoa(port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
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
}
