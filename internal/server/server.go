package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gertanoh/loadbalancer/internal/discovery"
)

type Config struct {
	HttpAddr string
	Handler  http.Handler
	discovery.Config
	MembershipHandler discovery.Handler
}

type Server struct {
	config     Config
	membership *discovery.Membership
}

func New(c Config) (*Server, error) {

	m, err := discovery.New(c.MembershipHandler, discovery.Config{
		BindAddr:       c.BindAddr,
		NodeName:       c.NodeName,
		Tags:           c.Tags,
		StartJoinAddrs: c.StartJoinAddrs,
	})

	if err != nil {
		return nil, err
	}

	srv := &Server{
		config:     c,
		membership: m,
	}
	return srv, nil
}

func (s *Server) Serve() error {
	srv := http.Server{
		Addr:         s.config.HttpAddr,
		Handler:      s.config.Handler,
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

		err = s.membership.Leave()
		if err != nil {
			shutdownError <- err
		}
		shutdownError <- nil
	}()

	// Start the HTTP server
	log.Println("Waiting for incoming connections at ", s.config.HttpAddr)
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
