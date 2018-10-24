package peek

import (
	"errors"
	"io"
	"net"
	"syscall"
)

var ErrNotSupported = errors.New("Operation not supported")

type SyscallConn interface {
	SyscallConn() (syscall.RawConn, error)
}
