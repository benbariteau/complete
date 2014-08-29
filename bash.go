package complete

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var funcs map[string]string

func init() {
	funcs = bashCompletionFuncs()
}

func Bash(line string) (completions []string) {
	fields := strings.Fields(line)
	bashCompletionFunc := funcs[fields[0]]

	completionsCmd := exec.Command(
		"env",
		fmt.Sprintf("COMP_CWORD=%v", len(strings.Fields(line))-1),
		fmt.Sprintf("COMP_LINE=%v", line),
		fmt.Sprintf("COMP_POINT=%v", len(line)+1),
		"bash", "-c",
		fmt.Sprint(
			fmt.Sprintf("%v;", print_completions_src),
			fmt.Sprintf("COMP_WORDS=(%v);", line),
			fmt.Sprintf(
				"source %v; %v; ",
				bashCompletionPath(),
				bashCompletionFunc,
			),
			"__print_completions;",
		),
	)

	stdout, err := completionsCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	completionsCmd.Start()

	out := bufio.NewReader(stdout)

	prefix := strings.Join(fields[:len(fields)-1], " ")
	for line, err := out.ReadString('\n'); err == nil; line, err = out.ReadString('\n') {
		completions = append(
			completions,
			fmt.Sprintf(
				"%v %v",
				prefix,
				line[:len(line)-1],
			),
		)
	}
	return
}

func bashCompletionFuncs() map[string]string {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %v; complete -p", bashCompletionPath()))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(stdout)

	cmd.Start()

	completionFuncs := make(map[string]string)

	rx := regexp.MustCompile("-F ([^ ]*).* ([^ ]*)$")
	for scanner.Scan() {
		match := rx.FindStringSubmatch(scanner.Text())
		if match != nil {
			funcName := match[1]
			cmdName := match[2]
			completionFuncs[cmdName] = funcName
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return completionFuncs
}

const bashCompletionPathUnix = "/etc/bash_completion"

func bashCompletionPath() string {
	switch runtime.GOOS {
	case "darwin":
		return brewPrefix() + bashCompletionPathUnix
	default:
		return bashCompletionPathUnix
	}
}

func brewPrefix() string {
	cmd := exec.Command("brew", "--prefix")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}

	cmd.Start()

	stdout := bufio.NewReader(stdoutPipe)
	line, err := stdout.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return ""
	}

	return line[:len(line)-1]
}

const print_completions_src = `__print_completions() { for ((i=0;i<${#COMPREPLY[*]};i++)); do echo ${COMPREPLY[i]};done; }`
