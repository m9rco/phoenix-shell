package shell

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strings"
)

type editor interface {
  ReadCode() (string, error)
}

type minEditor struct {
  in  *bufio.Reader
  out io.Writer
}

func newMinEditor(in, out *os.File) *minEditor {
  return &minEditor{bufio.NewReader(in), out}
}

func (ed *minEditor) ReadCode() (string, error) {
  wd, err := os.Getwd()
  if err != nil {
    wd = "?"
  }
  fmt.Fprintf(ed.out, "%s> ", wd)
  line, err := ed.in.ReadString('\n')
  line = strings.TrimSpace(line)
  line = strings.TrimRight(line, "\r\n\t")
  line = strings.TrimLeft(line, "\r\n\t")
  return line, err
}
