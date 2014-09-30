package complete

import (
	"os"
	"os/exec"
	"strings"
)

func Path(prefix string) (completions []string) {
	cmd := exec.Command("sh", "-c", "echo $PATH")
	outBytes, err := cmd.Output()
	if err != nil {
		return
	}
	out := string(outBytes)
	paths := strings.Split(out, ":")
	for _, path := range paths {
		dir, err := os.Open(path)
		if err != nil {
			continue
		}
		programs, err := dir.Readdirnames(0)
		if err != nil {
			continue
		}
		for _, program := range programs {
			if strings.HasPrefix(program, prefix) {
				completions = append(completions, program)
			}
		}
	}
	return
}
