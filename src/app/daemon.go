package app

import (
  "os"

  "github.com/m9rco/phoenix-shell/src/app/daemon"
  daemonsvc "github.com/m9rco/phoenix-shell/src/pkg/daemon"
)

type daemonProgram struct{ inner *daemon.Daemon }

func (p daemonProgram) Main(fds [3]*os.File, _ []string) int {
  err := p.inner.Main(daemonsvc.Serve)
  if err != nil {
    logger.Println("daemon error:", err)
    return 2
  }
  return 0
}
