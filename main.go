package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
)

var (
	version = "unknown"
	command []string
	program = filepath.Base(os.Args[0])
)

func parseFlags() {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	}

	command = os.Args[1:]
}

func showUsage() {
	fmt.Printf("Usage: %s <command> [args]...\n", program)
	fmt.Printf("Version: %s\n", version)
	fmt.Println("")
	fmt.Println("  Chronic runs the <command> and hides the output unless the command returns a non-zero exit code.")
}

func tempFile(prefix string) *os.File {
	var tempFile *os.File
	var err error

	if tempFile, err = ioutil.TempFile("", program+"-"+prefix); err != nil {
		fatal(err)
	}

	return tempFile
}

func emitCommand() {
	fmt.Printf("**** command ****\n")
	fmt.Printf("%#q\n", command)
	fmt.Println()
}

func emitOutput(name string, file io.ReadSeeker) {
	shownHeader := false

	if _, err := file.Seek(0, 0); err != nil {
		fatal(err)
	}
	buff := bufio.NewScanner(file)

	for buff.Scan() {
		if !shownHeader {
			fmt.Printf("**** %s ****\n", name)
			shownHeader = true
		}
		fmt.Printf("%s: %s\n", name, buff.Text())
	}

	if shownHeader {
		fmt.Println()
	}
}

func fatal(err error) int {
	fmt.Printf("[FATAL %s] %s\n", program, err)
	if user, err := user.Current(); err == nil {
		fmt.Printf("[FATAL %s] User:  %q (%s)\n", program, user.Username, user.Uid)
	}
	fmt.Printf("[FATAL %s] $PATH: %s\n", program, os.Getenv("PATH"))
	fmt.Printf("\n")
	showUsage()
	return -1
}

func run() int {
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var err error

	cmd := exec.Command(command[0], command[1:]...)
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
		if exiterr, ok := err.(*exec.ExitError); ok {
			emitCommand()
			emitOutput("stdout", tmpOut)
			emitOutput("stderr", tmpErr)

			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				ec := status.ExitStatus()
				fmt.Printf("Exited with %d\n", ec)
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
