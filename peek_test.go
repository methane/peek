package peek

import (
	"io"
	"net"
	"testing"
	"time"
)

// https://gist.github.com/tsavola/cd847385989f1ae497dbbcd2bba68753
func connPair() (serverConn, clientConn net.Conn, err error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return
	}
	defer l.Close()

	addr := l.Addr()

	var err2 error

	done := make(chan struct{})

	go func() {
		defer close(done)
		clientConn, err2 = net.Dial(addr.Network(), addr.String())
	}()

	serverConn, err = l.Accept()
	<-done

	if err == nil {
		err = err2
	}

	if err != nil {
		if clientConn != nil {
			clientConn.Close()
		}
		if serverConn != nil {
			serverConn.Close()
		}
	}

	return
}

func TestPeek(t *testing.T) {
	sc, cc, err := connPair()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()
	defer cc.Close()

	n, err := cc.Write([]byte("Hello"))
	if err != nil || n != 5 {
		t.Fatalf("Failed to send: n=%v, err=%v", n, err)
	}

	buf := make([]byte, 10)

	timedout := true
	for i := 0; i < 10; i++ {
		n, err = Peek(sc, buf)
		if err == nil && n == 0 {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		timedout = false
		break
	}
	if timedout {
		t.Fatal("timed out")
	}

	if err != nil || n != 5 {
		t.Fatalf("Failed to peek: n=%v, err=%v", n, err)
	}

	n, err = Peek(sc, buf)
	if err != nil || n != 0 {
		t.Fatalf("Failed to non-blocking read: n=%v, err=%v", n, err)
	}

	cc.Close()

	timedout = true
	for i := 0; i < 10; i++ {
		n, err = Peek(sc, buf)
		if err == nil && n == 0 {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		timedout = false
		break
	}
	if timedout {
		t.Fatal("timed out")
	}

	if err != io.EOF {
		t.Errorf("Unexpected err: %v", err)
	}
}
