package monitor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/joycastle/casual-server-lib/log"
)

func Exec(cmd string) ([]string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	if err := execCmd.Start(); err != nil {
		return nil, err
	}

	err := execCmd.Wait()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if stderrStr != "" {
		log.Get("error").Info("Exec", cmd, stderrStr)
	}

	ret := []string{}
	for _, v := range strings.Split(stdoutStr, "\n") {
		if len(v) > 0 {
			v = strings.Trim(v, " ")
			v = strings.Trim(v, "\t")
			v = strings.Trim(v, "\n")
			v = strings.Trim(v, "\r")
			ret = append(ret, v)
		}
	}

	return ret, err
}

func MergeCmd(cmd1, cmd2 string) string {
	return fmt.Sprintf("%s '%s'", cmd1, cmd2)
}
