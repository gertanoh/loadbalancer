package main

import (
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

// TODO user serf for services discovery here, now hard coded
var (
	servers = []string{"http://localhost:4001", "http://localhost:4002", "http://localhost:4003"}
	counter int32 // round robin index
)

func roundRobinSelector() string {
	value := atomic.AddInt32(&counter, 1)
	return servers[int(value)%len(servers)]
}
func main() {

	http.HandleFunc("/", handleRequestAndRedirect)
	log.Println("Load Balancer on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	log.Println("---------------------------")
	log.Println("Received request from ", r.RemoteAddr)
	log.Println(r.Method, r.URL.Path, r.Proto)

	for name, headers := range r.Header {
		for _, h := range headers {
			log.Println(name, h)
		}
	}

	// Forward
	// TODO round robin and servers discovery
	server := roundRobinSelector()
	log.Println("server addr: ", server)
	resp, err := http.Get(server + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(body)
}
