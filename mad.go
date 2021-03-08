package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// Enviroment contain the mad specific enviroment variables used
// in the execution of the current script
type Enviroment struct {
	IsPreview bool
	IsBlock   bool
	FullInput string
	InputLen  int
}

func (e Enviroment) Env() []string {
	return []string{

		fmt.Sprintf("MAD_ISPREVIEW=%t", e.IsPreview),
		fmt.Sprintf("MAD_ISBLOCK=%t", e.IsBlock),
		fmt.Sprintf("MAD_FULLINPUT=%s", e.FullInput),
		fmt.Sprintf("MAD_INPUTLEN=%d", e.InputLen),
	}
}

func main() {
	preview := flag.Bool("preview", false, "Render preview")
	flag.Parse()
	fl := flag.Arg(0)
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
	match := MakeMatch(madPath, *preview)
	out := parseReg.ReplaceAllStringFunc(contentString, match)
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

var parseReg = regexp.MustCompile(`(?ms)\[(.*)]:#\s*\(\s*(.*)\s*\)`)

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
func MakeMatch(path string, isPreview bool) func(string) string {
	return func(match string) string {
		cmd, arg := Parse(match)
		exe, found := SearchInPath(path, cmd)
		_ = arg
		if !found {
			fmt.Fprintln(os.Stderr, "Cannot find", cmd)
			return match
		}

		args := strings.Split(arg, " ")

		env := Enviroment{
			IsPreview: isPreview,
			IsBlock:   strings.ContainsRune(match, '\n'),
			FullInput: arg,
			InputLen:  len(args),
		}

		out, ok := Execute(exe, env, args...)
		if !ok {
			fmt.Fprintf(os.Stderr, "Cannot execute %s\n", exe)
		}

		if isPreview {
		}

		switch {
		case isPreview && out != "":
			return fmt.Sprintf("%s\n<!--\f%s\n-->", match, out)
		case isPreview && out == "":
			return match
		default:
			return out
		}
	}
}

//Execute executes the given command exe with the argument arg and returns the output
// if the command executes successfully ok is true and output contains stdout
// if the command doesn't execute successfully ok is false and output contains stderr
func Execute(exe string, env Enviroment, args ...string) (output string, ok bool) {
	c := exec.Command(exe, args...)
	c.Env = append(c.Env, env.Env()...)
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
