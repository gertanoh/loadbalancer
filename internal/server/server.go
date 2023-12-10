package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/serf/serf"
)

// Serf membership config
type Membership struct {
	NodeName       string            // node unique name
	BindAddr       string            // addr for gossiping
	Tag            map[string]string // use to share information, we will share each server addr for the load balancer here
	StartJoinAddrs []string

	handler SerfHandler
	serf    *serf.Serf
	events  chan serf.Event
}

type Config struct {
	HttpAddr     string
	Handler      http.Handler
	SerfNodeName string
	Membership
}

type Server struct {
	Config
}

// Serf service discovery

func NewSerfNode(handler SerfHandler, config Config) (*Membership, error) {
	c := &Membership{}
}

func (m *Membership) setupSerf() error {
	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)
	if err != nil {
		return err
	}

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	config.EventCh
}

func New(c Config) (*Server, error) {

	srv := &Server{
		Config: c,
	}

	return srv, nil
}

func (s *Server) Serve() error {
	srv := http.Server{
		Addr:         s.HttpAddr,
		Handler:      s.Handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("shutting down server")

		// create a 5s context to terminate the server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the server
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		shutdownError <- nil
	}()

	// Start the HTTP server
	log.Println("Waiting for incoming connections at ", s.HttpAddr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	err = <-shutdownError
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
