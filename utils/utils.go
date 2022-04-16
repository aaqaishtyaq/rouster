package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func RunCommand(command string, args []string) error {
	var a []string
	for _, str := range args {
		if str != "" {
			a = append(a, str)
		}
	}

	fmt.Println(a)
	cmd := exec.Command(command, a...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for command", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for command", err)
		return err
	}

	cmdReader := io.MultiReader(stderr, stdout)

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	log.Println("Started running shell command")

	fmt.Println("DEBUG-1")
	fmt.Println(err)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting command", err)
		return err
	}

	var waitStatus syscall.WaitStatus
	err = cmd.Wait()

	if exitError, ok := err.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
		if waitStatus.ExitStatus() != 0 {
			fmt.Println("DEBUG")
			fmt.Println(err)
			fmt.Fprintln(os.Stderr, "Error waiting command", err)
			fmt.Println("Error ", err.Error())
			return err
		}
	}

	return nil
}
