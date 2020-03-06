package shell

import (
  "fmt"
  "github.com/m9rco/phoenix-shell/src/pkg/sys"
  "golang.org/x/sys/unix"
  "io"
  "os"
  "os/exec"
  "os/user"
  "strings"
  "syscall"
  "time"
)

func interact(fds [3]*os.File) {
  var ed editor
  ed = newMinEditor(fds[0], fds[2])
  sanitize(fds[0], fds[2])
  cooldown := time.Second
  for {
    line, err := ed.ReadCode()
    if line == "" {
      continue
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
    if line == "exit" {
      break
    }

    if len(line) > 0 {
      _ = runCommand(line)
    }
    cooldown = time.Second
    sanitize(fds[0], fds[2])
  }
}

func switchDir(cmd []string) {
  var dir string
  if len(cmd) > 2 {
    fmt.Println("Too many arguments to builtin 'cd'.\n")
  } else if len(cmd) == 2 {
    dir = cmd[1]
  } else {
    u, err := user.Current()
    if err != nil {
      fmt.Println(fmt.Sprintf("Unable to get current user: %s\n", err))
    }
    dir = u.HomeDir
  }

  if dir == "~" {
    u, err := user.Current()
    if err != nil {
      fmt.Println(fmt.Sprintf("Unable to get current user: %s\n", err))
    }
    dir = u.HomeDir
  }

  if _, err := os.Stat(dir); err != nil {
    if os.IsNotExist(err) {
      wd, err := os.Getwd()
      if err != nil {
        wd = "/"
      }
      dir = wd + dir
      if _, err := os.Stat(dir); err != nil {
        fmt.Println("Directory does not exist")
        return
      }
    }
  }

  if err := os.Chdir(dir); err != nil {
    fmt.Println(fmt.Sprintf("Unable to change directory to '%s': %s\n", dir, err))
  }
  return
}

func runCommand(cmd string) (retval int) {
  if len(cmd) <= 0 {
    return
  }
  retval = sys.EXIT_SUCCESS
  _, err := os.Getwd()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to get current working directory: %v\n", err)
  }
  cmds := strings.Fields(cmd)
  given := strings.Join(cmds, " ")
  if len(given) < 1 {
    return
  }

  if cmds[0] == "cd" {
    switchDir(cmds)
    return
  }

  c := exec.Command(cmds[0], cmds[1:]...)
  c.Stdin = os.Stdin
  c.Stdout = os.Stdout
  c.Stderr = os.Stderr
  if err := c.Run(); err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      retval = exitError.Sys().(syscall.WaitStatus).ExitStatus()
    }
  }

  return
}

func sanitize(in, out *os.File) {
  _ = unix.SetNonblock(int(in.Fd()), false)
  _ = unix.SetNonblock(int(out.Fd()), false)
}
