// +build !windows,!plan9

package shell

import (
  "fmt"
  "github.com/m9rco/phoenix-shell/src/pkg/sys"
  "os"
  "syscall"
)

func handleSignal(sig os.Signal, stderr *os.File) {
  switch sig {
  case syscall.SIGHUP:
    _ = syscall.Kill(0, syscall.SIGHUP)
    os.Exit(0)
  case syscall.SIGUSR1:
    fmt.Fprint(stderr, sys.DumpStack())
  }
}
