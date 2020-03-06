package daemon

import (
  "os"
  "syscall"

  "golang.org/x/sys/unix"
)

func setUmask() {
  unix.Umask(0077)
}

func procAttrForSpawn() *os.ProcAttr {
  return &os.ProcAttr{
    Dir:   "/",
    Env:   []string{},
    Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
    Sys: &syscall.SysProcAttr{
      Setsid: true,
    },
  }
}
