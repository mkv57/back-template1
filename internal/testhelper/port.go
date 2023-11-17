package testhelper

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Host is std test host.
const Host = `127.0.0.1`

// Addr returns free test host:port.
func Addr(t *testing.T, assert *require.Assertions) netip.AddrPort {
	t.Helper()

	return netip.AddrPortFrom(netip.MustParseAddr(Host), UnusedTCPPort(t, assert, Host))
}

//nolint:gochecknoglobals // By design.
var (
	usedTCPPort   = make(map[int]int)
	usedTCPPortMu sync.Mutex
)

// UnusedTCPPort returns random unique unused TCP port at host.
func UnusedTCPPort(t *testing.T, assert *require.Assertions, host string) uint16 {
	t.Helper()

	var port int
	var portStr string
	ln, err := net.Listen("tcp", host+":0")
	if err == nil {
		err = ln.Close()
	}
	if err == nil {
		_, portStr, err = net.SplitHostPort(ln.Addr().String())
	}
	if err == nil {
		port, err = strconv.Atoi(portStr)
	}
	assert.NoError(err)

	usedTCPPortMu.Lock()
	used := usedTCPPort[port]
	usedTCPPort[port]++
	usedTCPPortMu.Unlock()
	if used > 0 {
		const maxRecursion = 3
		if used > maxRecursion {
			panic(fmt.Sprintf("same TCP port returned multiple times: %d", port))
		}

		return UnusedTCPPort(t, assert, host)
	}

	return uint16(port)
}

// WaitTCPPort tries to connect to addr until success or ctx.Done.
func WaitTCPPort(ctx context.Context, addr string) error {
	const delay = time.Second / 20
	var dialer net.Dialer
	for ; ctx.Err() == nil; time.Sleep(delay) {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			return conn.Close()
		}
	}

	return ctx.Err()
}
