package main

import (
  "github.com/m9rco/phoenix-shell/src/app"
  "os"
)

func main() {
  os.Exit(app.Main([3]*os.File{os.Stdin, os.Stdout, os.Stderr}, os.Args))
}
