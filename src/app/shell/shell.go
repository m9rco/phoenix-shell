package shell

import (
  "github.com/m9rco/phoenix-shell/src/pkg/sys"
  "github.com/m9rco/phoenix-shell/src/pkg/util"
  "os"
  "os/signal"
  "syscall"
)

var logger = util.GetLogger("[shell] ")

type Shell struct {
  BinPath     string
  SockPath    string
  DbPath      string
  Cmd         bool
  CompileOnly bool
  NoRc        bool
  JSON        bool
}

func (sh *Shell) Main(fds [3]*os.File, args []string) int {
  defer rescue()
  //restoreTTY := term.SetupGlobal()
  //defer restoreTTY()
  handleSignals(fds[2])
  interact(fds)
  return 0
}

func rescue() {
  if r := recover(); r != nil {
    println(r)
    print(sys.DumpStack())
    currentShell := os.Getenv("SHELL")
    if len(currentShell) <= 0 {
      currentShell = "/bin/sh"
    }
    println("\nexecing recovery shell " + currentShell)
    _ = syscall.Exec(currentShell, []string{currentShell}, os.Environ())
  }
}

func handleSignals(stderr *os.File) {
  sigs := make(chan os.Signal)
  signal.Notify(sigs)
  go func() {
    for sig := range sigs {
      logger.Println("signal", sig)
      handleSignal(sig, stderr)
    }
  }()
}
