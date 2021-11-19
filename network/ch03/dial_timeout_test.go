package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func DialTimeout(network, address string, timeout time.Duration,
) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "!connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	c, err := DialTimeout("tcp", "127.0.0.1:http", 2*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}
	t.Log(err)
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}
	t.Log("B")
	if !nErr.Timeout() {
		t.Fatal("error is not a timeout")
	}
	t.Log("C")

	time.Sleep(4 * time.Second)
}
