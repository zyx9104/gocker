package cmd

import (
	"fmt"
	"gocker/log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use: "exec [container id] [cmd]",
		Run: func(cmd *cobra.Command, args []string) {
			runExec(args[0], strings.Join(args[1:], " "))
		},
	}
)

func runExec(containerId, execCmd string) {
	c := e.GetContainer(containerId)
	if c == nil {
		log.Fatal("container not found")
	}
	cmd := exec.Command(selfProc, "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	os.Setenv("CONTAINER_PID", fmt.Sprintf("%v", c.Pid))
	os.Setenv("CONTAINER_CMD", execCmd)
	os.Setenv("PS1", fmt.Sprintf(`[container: %v]:\w \$ `, c.Id))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
