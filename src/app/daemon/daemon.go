package daemon

import (
  "errors"
  "fmt"
  "os"
  "path/filepath"
)
type Daemon struct {
  BinPath string
  DbPath string
  SockPath string
  LogPathPrefix string
}

func (d *Daemon) Main(serve func(string, string)) error {
  setUmask()
  serve(d.SockPath, d.DbPath)
  return nil
}

func (d *Daemon) Spawn() error {
  binPath := d.BinPath
  if binPath == "" {
    bin, err := os.Executable()
    if err != nil {
      return errors.New("cannot find phoenix-shell: " + err.Error())
    }
    binPath = bin
  }

  var pathError error
  abs := func(name string, path string) string {
    if pathError != nil {
      return ""
    }
    if path == "" {
      pathError = fmt.Errorf("%s is required for spawning daemon", name)
      return ""
    }
    absPath, err := filepath.Abs(path)
    if err != nil {
      pathError = fmt.Errorf("cannot resolve %s to absolute path: %s", name, err)
    }
    return absPath
  }
  binPath = abs("BinPath", binPath)
  dbPath := abs("DbPath", d.DbPath)
  sockPath := abs("SockPath", d.SockPath)
  logPathPrefix := abs("LogPathPrefix", d.LogPathPrefix)
  if pathError != nil {
    return pathError
  }

  args := []string{
    binPath,
    "-daemon",
    "-bin", binPath,
    "-db", dbPath,
    "-sock", sockPath,
    "-logprefix", logPathPrefix,
  }

  _, err := os.StartProcess(binPath, args, procAttrForSpawn())
  return err
}
