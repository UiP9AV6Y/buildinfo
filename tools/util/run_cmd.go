package util

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

// RunCmd is a wrapper around [exec.Command] and returns
// its standard output. Any error is the result of an
// execution failure or an error instance comprisef of
// the error output, should no standard output be produced.
func RunCmd(cmd string, argv []string) (string, error) {
	run := exec.Command(cmd, argv...)
	if run.Err != nil {
		return "", run.Err
	}

	stderr, err := run.StderrPipe()
	if err != nil {
		return "", err
	}

	stdout, err := run.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := run.Start(); err != nil {
		return "", err
	}

	e, err := io.ReadAll(stderr)
	if err != nil {
		return "", err
	}

	o, _ := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	outData := strings.Trim(string(o), " \n\r")
	errData := strings.Trim(string(e), " \n\r")
	err = run.Wait()
	if errData != "" && (err != nil || outData == "") {
		return "", errors.New(errData)
	} else if err != nil {
		return "", err
	}

	return outData, nil
}
