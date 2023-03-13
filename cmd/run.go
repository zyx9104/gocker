//go:build linux

package cmd

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"gocker/container"
	"gocker/log"
)

var (
	d      bool
	volume string
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			runRun(args)
		},
	}
)

func init() {
	runCmd.PersistentFlags().BoolVarP(&d, "d", "d", false, "后台运行")
	runCmd.PersistentFlags().StringVarP(&volume, "volume", "v", "", "volume")
}

func runRun(args []string) {
	if os.Args[0] == selfProc {
		run(args)
		return
	}
	c := e.RunContainer(args[0])
	args[0] = path.Join(container.BaseContainers, c.Id)
	cmd := exec.Command(selfProc, append([]string{"run"}, args...)...)
	log.Infof("-d %v", d)
	log.Infof("-v %v", volume)
	if !d {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWNS |
			unix.CLONE_NEWUTS |
			unix.CLONE_NEWIPC |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	c.Pid = cmd.Process.Pid
	log.Infof("parent pid: %v, container pid: %v", os.Getpid(), cmd.Process.Pid)

	// cgroup.MemoryTasks("/gocker/test", cmd.Process.Pid)
	// cgroup.Memory("/gocker/test", "100m")
	if volume != "" {
		ss := strings.Split(volume, ":")
		cp := path.Join(container.BaseContainers, c.Id, ss[1])
		unix.Mount(ss[0], cp, "", unix.MS_BIND, "")
	}
	if !d {
		cmd.Process.Wait()
	}
	log.Info("parent Done")
}

func run(args []string) {

	log.Infof("[os.Args] %v", os.Args)
	pivotRoot(args[0])

	abcCmd, err := exec.LookPath(args[1])
	args[1] = abcCmd
	if err != nil {
		log.Panic(err)
	}
	unix.Sethostname([]byte(uuid.NewString()))
	log.Infof("[CMD] %v", args[1:])
	if err := unix.Exec(abcCmd, args[1:], os.Environ()); err != nil {
		log.Fatal(err)
	}
}
