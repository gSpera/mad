package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	preview := flag.Bool("preview", false, "Render preview")
	debug := flag.Bool("debug", false, "Debug messages")
	flag.Parse()
	fl := flag.Arg(0)
	if !*debug {
		log.Default().SetOutput(io.Discard)
	}

	madPath := os.Getenv("MAD_PATH")
	if madPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get user home: MAD_PATH not defined\n")
			os.Exit(1)
		}

		madPath = path.Join(home, ".config", "mad", "bin")
	}

	env := Enviroment{
		IsPreview: *preview,

		searchPath: madPath,
	}

	content, err := ioutil.ReadFile(fl)
	if err != nil {
		panic(err)
	}
	contentRune := []rune(string(content))

	var out io.Writer = os.Stdout
	found := true
	var index, newIndex, nextIndex int
	var cmd Command

	log.Println("Init")
	for found {
		log.Println("Searching")
		newIndex, found = SearchCommand(contentRune, index)
		if !found {
			log.Println("\tNot found:", newIndex, found)

			// No new command, exit from loop and write everything
			log.Printf("Write: [%d:%d]: \"%s\"\n", index, len(contentRune), string(contentRune[index:]))
			writeRunes(out, contentRune[index:])
			break
		}
		log.Printf(" - %d:%d %t '%c'\n", index, newIndex, found, contentRune[newIndex])

		log.Println("Parsing")
		cmd, nextIndex, err = ParseCommand(contentRune, newIndex)
		if err != nil {
			log.Println("Cannot parse", err)
			log.Printf("Write: [%d:%d]: \"%s\"\n", index, nextIndex, string(contentRune[index:nextIndex]))
			// Not a command, write everything and the command and continue
			writeRunes(out, contentRune[index:nextIndex])

			index = nextIndex
			continue
		}
		log.Println(" ", cmd, newIndex, err)

		env.IsBlock = strings.ContainsRune(cmd.Arg, '\n')
		env.FullInput = cmd.Arg

		// Write anything
		writeRunes(out, contentRune[index:newIndex])
		// Execute
		ExecuteCommand(out, cmd, env)
		index = nextIndex
	}

	log.Println("Done")
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

// ExecuteCommand elaborate the match with the execution of the command
func ExecuteCommand(w io.Writer, cmd Command, env Enviroment) bool {
	exe, found := SearchInPath(env.searchPath, cmd.ScriptName)
	if !found {
		// Write source
		log.Println("Not found script")
		fmt.Fprint(w, cmd.Source)
		return false
	}

	args := strings.Split(cmd.Arg, " ")

	out, ok := Execute(exe, args, env)
	if !ok {
		fmt.Fprintf(os.Stderr, "Cannot execute %s\n", exe)
		fmt.Fprint(w, cmd.Source)
		return false
	}

	switch {
	case env.IsPreview && out != "":
		fmt.Fprintf(w, "%s\n<!--\n%s\n-->", cmd.Source, out)
	case env.IsPreview && out == "":
		fmt.Fprint(w, cmd.Source)
	default:
		fmt.Fprint(w, out)
	}

	return true
}

//Execute executes the given command exe with the argument arg and returns the output
// if the command executes successfully ok is true and output contains stdout
// if the command doesn't execute successfully ok is false and output contains stderr
func Execute(exe string, args []string, env Enviroment) (output string, ok bool) {
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
	case errors.Is(err, &exec.ExitError{}):
		exit := new(exec.ExitError)
		errors.As(err, &exit)
		return stderr.String(), false
	default:
		fmt.Fprintf(os.Stderr, "Cannot execute: error: %s\n%s\n", err, stderr.String())
		return stderr.String(), false
	}
}

// SearchCommand searchs for a plausible command in command, if not found returns the end and false
func SearchCommand(command []rune, startIndex int) (foundIndex int, ok bool) {
	if startIndex >= len(command) {
		return startIndex, false
	}

	index := startIndex

	for command[index] != '[' {
		index++

		if index >= len(command) {
			return index, false
		}
	}

	return index, true
}

func writeRunes(w io.Writer, c []rune) error {
	content := []byte(string(c))

	_, e := w.Write(content)
	return e
}
