package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gertanoh/loadbalancer/internal/helpers"
	"github.com/gertanoh/loadbalancer/internal/server"
	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
)

const (
	port = 8080
)

var srvPool *helpers.ServerPool

type SerfClusterHandler struct{}

func main() {

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Println("Load Balancer on port 8080")

	srvPool = helpers.NewServerPool()

	handler := handleConnection()

	var config server.Config
	config.HttpAddr = fmt.Sprintf(":%d", port)
	config.Handler = handler
	config.BindAddr = fmt.Sprintf("127.0.0.1:%d", 8400)
	config.NodeName = fmt.Sprintf("agent-%s", uuid.New())
	config.Tags = map[string]string{
		"http_addr":        config.HttpAddr,
		"is_load_balancer": "true",
	}
	serfHandler := &SerfClusterHandler{}
	config.MembershipHandler = serfHandler

	srv, err := server.New(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Serve(); err != nil {
		log.Fatal(err)
	}
}

func handleConnection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("---------------------------")
		log.Println("Received request from ", r.RemoteAddr)
		log.Println(r.Method, r.URL.Path, r.Proto)

		for name, headers := range r.Header {
			for _, h := range headers {
				log.Println(name, h)
			}
		}

		server, err := srvPool.GetNextServer()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		log.Println("server addr: ", server)
		address := "http://" + server.Address
		resp, err := http.Get(address + r.URL.Path)
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
	})
}

func (h *SerfClusterHandler) Join(id, addr, isLoadBalancer string) error {
	srv := helpers.Server{
		ID:      id,
		Address: addr,
	}
	srvPool.AddServer(srv)
	log.Println("Join cluster ", id, addr)
	return nil
}

func (h *SerfClusterHandler) Leave(id string) error {
	srvPool.RemoveServer(id)
	log.Println("left cluster")
	return nil
}
