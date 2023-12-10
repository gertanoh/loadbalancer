package helpers

import (
	"errors"
	"log"
	"net"
	"sync"
)

// getRandomPort starts a listener on a random port
// and returns the port number.
func GetRandomPort() int {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}

type Server struct {
	ID      string
	Address string
}

type ServerPool struct {
	servers []Server
	m       sync.RWMutex
	index   int
}

func NewServerPool() *ServerPool {
	return &ServerPool{}
}

func (sp *ServerPool) AddServer(server Server) {
	sp.m.Lock()
	defer sp.m.Unlock()
	sp.servers = append(sp.servers, server)
}

func (sp *ServerPool) GetNextServer() (Server, error) {
	sp.m.Lock()
	defer sp.m.Unlock()

	if len(sp.servers) == 0 {
		return Server{}, errors.New("empty servers pool")
	}

	server := sp.servers[sp.index]
	sp.index = (sp.index + 1) % len(sp.servers)
	return server, nil
}

func (sp *ServerPool) RemoveServer(serverID string) {
	sp.m.Lock()
	defer sp.m.Unlock()

	for i, server := range sp.servers {
		if server.ID == serverID {
			// remove server from slice
			sp.servers = append(sp.servers[:i], sp.servers[i+1:]...)
			// adjust index for round-robin
			if sp.index >= len(sp.servers) {
				sp.index = 0
			}
			return
		}
	}
}
