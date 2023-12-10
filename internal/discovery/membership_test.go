package discovery_test

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	. "github.com/gertanoh/loadbalancer/internal/discovery"
	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/require"
)

func TestMembership(t *testing.T) {
	m, handler := setupMember(t, nil, false)
	m, _ = setupMember(t, m, false)
	m, _ = setupMember(t, m, false)

	require.Eventually(t, func() bool {
		return len(handler.joins) == 2 &&
			len(m[0].Members()) == 3 &&
			len(handler.leaves) == 0
	}, 3*time.Second, 250*time.Millisecond)

	require.NoError(t, m[2].Leave())
	require.Eventually(t, func() bool {
		return len(handler.joins) == 2 &&
			len(m[0].Members()) == 3 &&
			serf.StatusLeft == m[0].Members()[2].Status &&
			len(handler.leaves) == 1
	}, 3*time.Second, 250*time.Millisecond)

	// Add the load balancer
	m, _ = setupMember(t, m, true)
	require.Eventually(t, func() bool {
		return len(handler.joins) == 3 &&
			len(m[0].Members()) == 4 &&
			m[0].Members()[3].Tags["is_load_balancer"] == "true" &&
			len(handler.leaves) == 1
	}, 3*time.Second, 250*time.Millisecond)

	require.Equal(t, fmt.Sprintf("%d", 2), <-handler.leaves)
}

// getRandomPort starts a listener on a random port
// and returns the port number.
func getRandomPort() int {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}

func setupMember(t *testing.T, members []*Membership, isLoadBalancer bool) ([]*Membership,
	*handler) {

	id := len(members)
	port := getRandomPort()
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	tags := map[string]string{
		"rpc_addr":         addr,
		"is_load_balancer": strconv.FormatBool(isLoadBalancer),
	}

	c := Config{
		NodeName: fmt.Sprintf("%d", id),
		BindAddr: addr,
		Tags:     tags,
	}

	h := &handler{}
	if len(members) == 0 {
		h.joins = make(chan map[string]string, 5)
		h.leaves = make(chan string, 5)
	} else {
		c.StartJoinAddrs = []string{
			members[0].BindAddr,
		}
	}
	m, err := New(h, c)
	require.NoError(t, err)
	members = append(members, m)
	return members, h
}

type handler struct {
	joins  chan map[string]string
	leaves chan string
}

func (h *handler) Join(id, addr, isLoadBalancer string) error {
	if h.joins != nil {
		h.joins <- map[string]string{
			"id":    id,
			"addr":  addr,
			"is_lb": isLoadBalancer,
		}
	}
	return nil
}

func (h *handler) Leave(id string) error {
	if h.leaves != nil {
		h.leaves <- id
	}
	return nil
}
