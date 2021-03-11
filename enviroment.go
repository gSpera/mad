package main

import "fmt"

// Enviroment contain the mad specific enviroment variables used
// in the execution of the current script
type Enviroment struct {
	IsPreview bool
	IsBlock   bool
	FullInput string
	InputLen  int

	// internal
	searchPath string
}

// Env returns a []string in the format accepted by exec.Env
// (for example MAD_ISPREVIEW=true)
func (e Enviroment) Env() []string {
	return []string{
		fmt.Sprintf("MAD_ISPREVIEW=%t", e.IsPreview),
		fmt.Sprintf("MAD_ISBLOCK=%t", e.IsBlock),
		fmt.Sprintf("MAD_FULLINPUT=%s", e.FullInput),
		fmt.Sprintf("MAD_INPUTLEN=%d", e.InputLen),
	}
}
