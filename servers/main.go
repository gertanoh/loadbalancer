package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gertanoh/loadbalancer/internal/helpers"
	"github.com/gertanoh/loadbalancer/internal/server"
	"github.com/google/uuid"
)

var port int
var host string
var clusterAddr string

type SerfClusterHandler struct{}

func main() {

	flag.IntVar(&port, "port", 8080, "API server port")
	flag.StringVar(&host, "host-IP", "127.0.0.1", "Host IP address")
	flag.StringVar(&clusterAddr, "cluster-addr", "", "Cluster Address")
	flag.Parse()

	if clusterAddr == "" {
		log.Fatal("Cluster address was not provided")
	}
	handler := handleRequests()
	var config server.Config
	config.HttpAddr = fmt.Sprintf("%s:%d", host, port)
	config.Handler = handler
	config.BindAddr = fmt.Sprintf("127.0.0.1:%d", helpers.GetRandomPort())
	config.NodeName = fmt.Sprintf("agent-%s", uuid.New())
	config.Tags = map[string]string{
		"http_addr":        config.HttpAddr,
		"is_load_balancer": "false",
	}
	serfHandler := &SerfClusterHandler{}
	config.MembershipHandler = serfHandler
	config.StartJoinAddrs = append(config.StartJoinAddrs, clusterAddr)

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

func (h *SerfClusterHandler) Join(id, addr, isLoadBalancer string) error {
	log.Println("Join cluster ", id, addr)
	return nil
}

func (h *SerfClusterHandler) Leave(id string) error {
	log.Println("left cluster ", id)
	return nil
}
