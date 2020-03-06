package daemon

import (
  "os"
  "syscall"
)

func setUmask() {
}

const (
  CREATE_BREAKAWAY_FROM_JOB = 0x01000000
  CREATE_NEW_PROCESS_GROUP  = 0x00000200
  DETACHED_PROCESS          = 0x00000008

  DaemonCreationFlags = CREATE_BREAKAWAY_FROM_JOB | CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS
)

func procAttrForSpawn() *os.ProcAttr {
  return &os.ProcAttr{
    Dir:   `C:\`,
    Env:   []string{"SystemRoot=" + os.Getenv("SystemRoot")},
    Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
    Sys:   &syscall.SysProcAttr{CreationFlags: DaemonCreationFlags},
  }
}
