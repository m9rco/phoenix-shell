package sys

import (
  "golang.org/x/sys/unix"
  "os"
  "unsafe"
)

// IsATTY returns true if the given file is a terminal.
func IsATTY(file *os.File) bool {
  type Termios unix.Termios
  var term Termios
  err := Ioctl(int(file.Fd()), unix.TIOCGETA, uintptr(unsafe.Pointer(&term)))
  return err == nil
}
