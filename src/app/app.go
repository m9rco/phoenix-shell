package app

import (
  "flag"
  "fmt"
  "github.com/m9rco/phoenix-shell/src/app/daemon"
  "github.com/m9rco/phoenix-shell/src/app/shell"
  "github.com/m9rco/phoenix-shell/src/pkg/util"
  "io"
  "log"
  "os"
  "runtime/pprof"
  "strconv"
)

const defaultWebPort = 3171
var logger = util.GetLogger("[main] ")

type flagSet struct {
  flag.FlagSet

  Log, LogPrefix, CPUProfile string

  Help, Version, BuildInfo, JSON bool

  CodeInArg, CompileOnly, NoRc bool

  Web  bool
  Port int

  Daemon bool
  Forked int

  Bin, DB, Sock string
}

func newFlagSet(stderr io.Writer) *flagSet {
  f := flagSet{}
  f.Init("elvish", flag.ContinueOnError)
  f.SetOutput(stderr)
  f.Usage = func() { usage(stderr, &f) }

  f.StringVar(&f.Log, "log", "", "a file to write debug log to")
  f.StringVar(&f.LogPrefix, "logprefix", "", "the prefix for the daemon log file")
  f.StringVar(&f.CPUProfile, "cpuprofile", "", "write cpu profile to file")

  f.BoolVar(&f.Help, "help", false, "show usage help and quit")
  f.BoolVar(&f.Version, "version", false, "show version and quit")
  f.BoolVar(&f.BuildInfo, "buildinfo", false, "show build info and quit")
  f.BoolVar(&f.JSON, "json", false, "show output in JSON. Useful with -buildinfo.")

  f.BoolVar(&f.CodeInArg, "c", false, "take first argument as code to execute")
  f.BoolVar(&f.CompileOnly, "compileonly", false, "Parse/Compile but do not execute")
  f.BoolVar(&f.NoRc, "norc", false, "run elvish without invoking rc.elv")

  f.BoolVar(&f.Web, "web", false, "run backend of web interface")
  f.IntVar(&f.Port, "port", defaultWebPort, "the port of the web backend")

  f.BoolVar(&f.Daemon, "daemon", false, "run daemon instead of shell")

  f.StringVar(&f.Bin, "bin", "", "path to the elvish binary")
  f.StringVar(&f.DB, "db", "", "path to the database")
  f.StringVar(&f.Sock, "sock", "", "path to the daemon socket")

  return &f
}

func Main(fds [3]*os.File, args []string) int {
  flag := newFlagSet(fds[2])
  err := flag.Parse(args[1:])
  if err != nil {
    return 2
  }
  if flag.CPUProfile != "" {
    f, err := os.Create(flag.CPUProfile)
    if err != nil {
      log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  if flag.Log != "" {
    err = util.SetOutputFile(flag.Log)
  } else if flag.LogPrefix != "" {
    err = util.SetOutputFile(flag.LogPrefix + strconv.Itoa(os.Getpid()))
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
  }

  return FindProgram(flag).Main(fds, flag.Args())
}

type Program interface {
  Main(fds [3]*os.File, args []string) int
}

// FindProgram finds a suitable Program according to flags. It does not have any
// side effects.
func FindProgram(flag *flagSet) Program {
  switch {
  case flag.Daemon:
    if len(flag.Args()) > 0 {
      return badUsageProgram{"arguments are not allowed with -daemon", flag}
    }
    return daemonProgram{&daemon.Daemon{
      BinPath:       flag.Bin,
      DbPath:        flag.DB,
      SockPath:      flag.Sock,
      LogPathPrefix: flag.LogPrefix,
    }}
  default:
    return &shell.Shell{
      BinPath: flag.Bin, SockPath: flag.Sock, DbPath: flag.DB,
      Cmd: flag.CodeInArg, CompileOnly: flag.CompileOnly,
      NoRc: flag.NoRc, JSON: flag.JSON}
  }
}
