package leader

import (
	"fmt"
	"syscall"
)

const path = "/tmp/%s-leader"

// Acquire attempts to place an exclusive lock on a temporary file
// shared by multiple services on the same machine.
// See https://man7.org/linux/man-pages/man2/flock.2.html for details.
func Acquire(service string) (int, error) {
	// Open the file with read permissions.
	fd, err := syscall.Open(fmt.Sprintf(path, service), syscall.O_CREAT|syscall.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}

	// Attempt to acquire a lock.
	if err := syscall.Flock(fd, syscall.LOCK_EX); err != nil {
		Release(fd)
		return 0, err
	}

	return fd, nil
}

// Release closes the underlying fd of the locked file and subsequently
// releases the lock.
func Release(fd int) {
	syscall.Close(fd)
}
