// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package peek

import (
	"syscall"
)

// Peek reads data from c into buff.
// Unlike Read, Peek doesn't wait receiving data (aka, non-blocking)
func Peek(c net.Conn, buff []byte) (n int, err error) {
	sc, ok := c.(SyscallConn)
	if !ok {
		return 0, ErrNotSupported
	}

	rc, err := sc.SyscallConn()
	if err != nil {
		return 0, ErrNotSupported
	}

	return peek(rc, buff)
}

func peek(rc syscall.RawConn, buff []byte) (n int, err error) {
	rerr := rc.Read(func(fd uintptr) bool {
		n, err = syscall.Read(int(fd), buff)
		return true
	})

	if rerr != nil {
		return n, rerr
	}
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	if err == syscall.EAGAIN {
		err = nil
		n = 0
	}
	return n, err
}
