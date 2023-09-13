package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

var (
	version = "unknown"
	command []string
	program = filepath.Base(os.Args[0])
	timeout = 0
)

func parseFlags() {
	flag.IntVar(&timeout, "t", 0, "Command timeout")
	flag.Parse()

	command = flag.Args()

	if len(command) == 0 {
		showUsage()
		os.Exit(0)
	}
}

func showUsage() {
	fmt.Fprintf(os.Stdout, "Usage: %s [-t TIMEOUT] <command> [args]...\n", program)
	fmt.Fprintf(os.Stdout, "Version: %s\n", version)
	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "  Chronic runs the <command> and hides the output unless the command returns a non-zero exit code. You can specify a timeout, after that, the process will be killed.\n")
}

func tempFile(prefix string) *os.File {
	var tempFile *os.File
	var err error

	if tempFile, err = os.CreateTemp("", program+"-"+prefix); err != nil {
		fatal(err)
	}

	return tempFile
}

func emitCommand(killed bool) {
	if killed {
		fmt.Fprintf(os.Stdout, "**** Timeout, process killed! ****\n\n")
	}
	fmt.Fprintf(os.Stdout, "**** command ****\n")
	fmt.Fprintf(os.Stdout, "%#q\n", command)
	fmt.Fprintf(os.Stdout, "\n")
}

func emitOutput(name string, file io.ReadSeeker) {
	shownHeader := false

	if _, err := file.Seek(0, 0); err != nil {
		fatal(err)
	}
	buff := bufio.NewScanner(file)

	for buff.Scan() {
		if !shownHeader {
			fmt.Fprintf(os.Stdout, "**** %s ****\n", name)
			shownHeader = true
		}
		fmt.Fprintf(os.Stdout, "%s: %s\n", name, buff.Text())
	}

	if shownHeader {
		fmt.Fprintf(os.Stdout, "\n")
	}
}

func fatal(err error) int {
	fmt.Fprintf(os.Stdout, "[FATAL %s] %s\n", program, err)
	if user, err := user.Current(); err == nil {
		fmt.Fprintf(os.Stdout, "[FATAL %s] User:  %q (%s)\n", program, user.Username, user.Uid)
	}
	fmt.Fprintf(os.Stdout, "[FATAL %s] $PATH: %s\n", program, os.Getenv("PATH"))
	fmt.Fprintf(os.Stdout, "\n")
	showUsage()

	return -1
}

func run() int {
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var err error

	cmd := exec.Command(command[0], command[1:]...) /* #nosec G204 */
	cmd.Stdin = os.Stdin

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return fatal(err)
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return fatal(err)
	}

	if err = cmd.Start(); err != nil {
		return fatal(err)
	}

	killed := false
	if timeout > 0 {
		var timer *time.Timer
		timer = time.AfterFunc(time.Duration(timeout)*time.Second, func() {
			timer.Stop()
			killed = true
			cmd.Process.Kill()
		})
	}

	tmpOut := tempFile("stdout")
	defer os.Remove(tmpOut.Name())
	if _, err = io.Copy(tmpOut, stdout); err != nil {
		return fatal(err)
	}

	tmpErr := tempFile("stderr")
	defer os.Remove(tmpErr.Name())
	if _, err = io.Copy(tmpErr, stderr); err != nil {
		return fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			emitCommand(killed)
			emitOutput("stdout", tmpOut)
			emitOutput("stderr", tmpErr)

			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				ec := status.ExitStatus()
				fmt.Fprintf(os.Stdout, "Exited with %d\n", ec)

				return ec
			}
		} else {
			return fatal(err)
		}
	}

	return 0
}

func main() {
	parseFlags()

	os.Exit(run())
}
