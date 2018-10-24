// +build windows

package peek

import (
	"syscall"
)

// Peek reads data from c into buff.
// Unlike Read, Peek doesn't wait receiving data (aka, non-blocking)
func Peek(c net.Conn, buff []byte) (n int, err error) {
	return 0, ErrNotSupported
}
