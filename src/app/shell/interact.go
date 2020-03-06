package shell

import (
  "fmt"
  "golang.org/x/sys/unix"
  "io"
  "os"
  "time"
)

func interact(fds [3]*os.File) {
  var ed editor
  ed = newMinEditor(fds[0], fds[2])
  sanitize(fds[0], fds[2])
  cooldown := time.Second
  cmdNum := 0

  for {
    cmdNum++
    line, err := ed.ReadCode()
    fmt.Println(line)
    if line == "exit" {
      break
    }
    if err == io.EOF {
      continue
    } else if err != nil {
      fmt.Fprintln(fds[2], "Editor error:", err)
      if _, isMinEditor := ed.(*minEditor); !isMinEditor {
        fmt.Fprintln(fds[2], "Falling back to basic line editor")
        ed = newMinEditor(fds[0], fds[2])
      } else {
        fmt.Fprintln(fds[2], "Don't know what to do, pid is", os.Getpid())
        fmt.Fprintln(fds[2], "Restarting editor in", cooldown)
        time.Sleep(cooldown)
        if cooldown < time.Minute {
          cooldown *= 2
        }
      }
      continue
    }

    cooldown = time.Second
    sanitize(fds[0], fds[2])
  }
}

func sanitize(in, out *os.File) {
  _ = unix.SetNonblock(int(in.Fd()), false)
  _ = unix.SetNonblock(int(out.Fd()), false)
}
