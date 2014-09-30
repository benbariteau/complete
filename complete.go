package complete

import (
	"strings"
)

/*
Complete returns a list of tab completions for a given partial command.

If line is one word (as defined by spaces), then a search of the PATH will be done, otherwise an attempt to call a bash completion function will be done.
*/
func Complete(line string) []string {
	parts := strings.Split(line, " ")
	if len(parts) == 1 && parts[0] != "" {
		return Path(parts[0])
	}

	return Bash(line)
}
