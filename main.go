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
	fmt.Printf("Usage: %s <command> [args]...\n\n", program)
	fmt.Println("Chronic runs the <command> and hides the output unless the command returns a non-zero exit code.")
}

func tempFile(prefix string) *os.File {
	var tempFile *os.File
	var err error

	if tempFile, err = ioutil.TempFile("", program+"-"+prefix); err != nil {
		fatal(err)
	}
	defer os.Remove(tempFile.Name())
	return tempFile
}

func emitCommand() {
	fmt.Printf("**** command ****\n")
	fmt.Printf("%#q\n", command)
	fmt.Println()
}

func emitOutput(name string, file *os.File) {
	shownHeader := false

	file.Seek(0, 0)
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

func fatal(err error) {
	fmt.Printf("[FATAL %s] %s\n", program, err)
	if user, err := user.Current(); err == nil {
		fmt.Printf("[FATAL %s] User:  %q (%s)\n", program, user.Username, user.Uid)
	}
	fmt.Printf("[FATAL %s] $PATH: %s\n", program, os.Getenv("PATH"))
	fmt.Printf("\n")
	showUsage()
	os.Exit(1)
}

func runCommand() int {
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var err error

	cmd := exec.Command(command[0], command[1:]...)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		fatal(err)
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		fatal(err)
	}

	if err = cmd.Start(); err != nil {
		fatal(err)
	}

	tmpOut := tempFile("stdout")
	io.Copy(tmpOut, stdout)

	tmpErr := tempFile("stderr")
	io.Copy(tmpErr, stderr)

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
			fatal(err)
		}
	}
	return 0
}

func main() {
	parseFlags()
	os.Exit(runCommand())
}
