package main

import (
	"fmt"
	"strings"
	"unicode"
)

// "[cmd] #? (TOKEN arg TOKEN)"
// cmd filename
// TOKEN ""
//		 "-X"
//		 "X-"
// arg

// Command rapresent a textual command
type Command struct {
	ScriptName string
	Arg        string
	IsBlock    bool

	Source string
}

// ParseCommand parses a full command
func ParseCommand(command []rune, startIndex int) (cmd Command, index int, err error) {
	index = startIndex
	ch := command[index]

	if ch != '[' {
		return cmd, index, fmt.Errorf("Found: '%c' expected: '['", ch)
	}
	index++

	// CMD
	cmdEnd, ok := ParseScriptName(command, index)
	if !ok {
		return cmd, index, fmt.Errorf("Cmd cannot be parsed")
	}

	cmd.ScriptName = string(command[index:cmdEnd])
	cmd.ScriptName = strings.TrimSpace(cmd.ScriptName)

	index = cmdEnd
	ch = command[index]

	// ]
	if ch != ']' {
		return cmd, index, fmt.Errorf("Found: '%c' expected: ']'", ch)
	}
	index++

	index = parseSkipSpace(command, index)
	ch = command[index]

	// #
	if ch == '#' {
		cmd.IsBlock = true
		index++
	}

	index = parseSkipSpace(command, index)
	ch = command[index]
	// (
	if ch != '(' {
		return cmd, index, fmt.Errorf("Found: '%c' expected: '('", ch)
	}
	index++

	// Arg
	argEnd, ok := ParseArg(command, index)
	if !ok {
		return cmd, index, fmt.Errorf("Argument cannot be parsed")
	}
	cmd.Arg = string(command[index:argEnd])
	cmd.Arg = strings.TrimSpace(cmd.Arg)
	index = argEnd
	ch = command[index]

	// )
	if ch != ')' {
		return cmd, index, fmt.Errorf("Found: '%c' expected: ')'", ch)
	}
	index++

	cmd.Source = string(command[startIndex:index])
	return cmd, index, nil
}

// ParseScriptName parses the command, returns the new index and a bool value
// if ok is false then ParseScriptName couldn't parse the command
func ParseScriptName(command []rune, startIndex int) (endIndex int, ok bool) {
	index := startIndex
	ch := command[index]

	for unicode.IsPrint(ch) && ch != ']' {
		index++
		if index >= len(command) {
			return 0, false
		}

		ch = command[index]
	}

	return index, true
}

// ParseArg parses the argument, new index and a bool value are returned
// if ok is false the ParseArg couldn't parse the command
func ParseArg(command []rune, startIndex int) (endIndex int, ok bool) {
	index := startIndex
	var level int
	ch := command[index]

	for unicode.IsPrint(ch) || unicode.IsSpace(ch) || ch == '\n' {
		if ch == '(' {
			level++
		} else if ch == ')' && level > 0 {
			level--
		} else if ch == ')' {
			break
		}

		index++
		if index >= len(command) {
			return 0, false
		}

		ch = command[index]
	}

	return index, true
}

func parseSkipSpace(command []rune, index int) int {
	for unicode.IsSpace(command[index]) {
		index++
	}

	return index
}
