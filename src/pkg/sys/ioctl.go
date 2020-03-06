package sys

import (
  "golang.org/x/sys/unix"
  "os"
)

func Ioctl(fd int, req uintptr, arg uintptr) error {
  _, _, e := unix.Syscall(
    unix.SYS_IOCTL, uintptr(fd), req, arg)
  if e != 0 {
    return os.NewSyscallError("ioctl", e)
  }
  return nil
}
