package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	if os.Args[0] == "/proc/self/exe" {
		//已经在容器内
		initContainer()
		unix.Exec(os.Args[1], os.Args[1:], os.Environ())
		return
	}
	cmd := exec.Command("/proc/self/exe", "/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{`PS1=[container] $ `}
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

	baseImage := "/home/zyx/gocker/tmp/rootfs"
	must(os.MkdirAll("/tmp/gocker/mnt", 0755))
	must(os.MkdirAll("/tmp/gocker/write", 0755))
	mountAufs("/tmp/gocker/mnt", "/tmp/gocker/write", []string{baseImage})

	cmd.Run()
}

func initContainer() {
	putOld := ".put_old"
	newRoot := "/tmp/gocker/mnt"
	oldRoot := path.Join(newRoot, putOld)
	unix.Mount(newRoot, newRoot, "", unix.MS_BIND|unix.MS_REC, "")
	unix.Mount("proc", path.Join(newRoot, "/proc"), "proc", 0, "")
	os.MkdirAll(oldRoot, 0755)
	must(unix.PivotRoot(newRoot, oldRoot))
	os.Chdir("/")
	unix.Unmount(putOld, unix.MNT_DETACH)
	os.RemoveAll(putOld)
}

func mountAufs(target, rwLayers string, layers []string) {
	roBranch := ""
	for _, s := range layers {
		roBranch += fmt.Sprintf("%v=ro:", s)
	}
	rw := fmt.Sprintf("%v=rw", rwLayers)
	branches := fmt.Sprintf("br:%v:%v", rw, roBranch)
	if err := unix.Mount("none", target, "aufs", 0, branches); err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
