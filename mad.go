package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func main() {
	fl := os.Args[1]
	madPath := os.Getenv("MAD_PATH")
	if madPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get user home: MAD_PATH not defined\n")
			os.Exit(1)
		}

		madPath = path.Join(home, ".config", "mad", "bin")
	}

	content, err := ioutil.ReadFile(fl)
	if err != nil {
		panic(err)
	}
	contentString := string(content)
	reg := regexp.MustCompile(`\[.*\]:# \(?(.*)\)?`)
	match := MakeMatch(madPath)
	out := reg.ReplaceAllStringFunc(contentString, match)
	fmt.Println(out)
}

// SearchInPath searches in the path the exe filename
// returns the path to the exe
func SearchInPath(searchPath string, exe string) (string, bool) {
	paths := strings.Split(searchPath, string(os.PathListSeparator))
	for _, dir := range paths {
		file := path.Join(dir, exe)
		_, err := os.Stat(file)

		// File doesn't exist
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		//Found
		if err == nil {
			return file, true
		}
	}

	// Not found
	return "", false
}

var parseReg = regexp.MustCompile(`\[(.*)]:# \((.*)\)`)

// Parse parses the splits the matched string [CMD]:# (ARG)
// in CMD and ARG
func Parse(input string) (cmd string, arg string) {
	ret := parseReg.FindStringSubmatch(input)
	cmd = ret[1]
	arg = ret[2]
	return
}

// MakeMatch builds Match(string) string
// which elaborate the match with the execution of the command
func MakeMatch(path string) func(string) string {
	return func(match string) string {
		cmd, arg := Parse(match)
		exe, found := SearchInPath(path, cmd)
		_ = arg
		if !found {
			fmt.Fprintln(os.Stderr, "Cannot find", cmd)
			return match
		}
		out, ok := Execute(exe, arg)
		_ = ok
		return out
	}
}

//Execute executes the given command exe with the argument arg and returns the output
// if the command executes successfully ok is true and output contains stdout
// if the command doesn't execute successfully ok is false and output contains stderr
func Execute(exe, arg string) (output string, ok bool) {
	args := strings.Split(arg, " ")
	args = append([]string{arg}, args...)

	c := exec.Command(exe, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c.Stdout = stdout
	c.Stderr = stderr

	err := c.Run()
	switch {
	case err == nil:
		return stdout.String(), true
	case errors.Is(err, &exec.ExitError{ProcessState: nil}):
		exit := new(exec.ExitError)
		errors.As(err, &exit)
		return stderr.String(), false
	default:
		fmt.Fprintf(os.Stderr, "Cannot execute: error: %s\n%s\n", err, stderr.String())
		return stderr.String(), false
	}
}
