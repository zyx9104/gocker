package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"gocker/engine"
	"gocker/log"
)

var (
	rootCmd = cobra.Command{
		Use: "gocker",
	}
	e = engine.New()
)

const (
	cgroupRoot = "/sys/fs/cgroup/"
	selfProc   = "/proc/self/exe"
)

func pivotRoot(newRoot string) {

	putOld := ".pivot_root"
	oldRoot := path.Join(newRoot, putOld)
	log.Infof("new_root: %v", newRoot)
	flags := uintptr(0)
	err := unix.Mount("proc", path.Join(newRoot, "/proc"), "proc", flags, "")
	if err != nil {
		log.Panic(err)
	}
	err = unix.Mount(newRoot, newRoot, "", unix.MS_BIND|unix.MS_REC, "")
	if err != nil {
		log.Panic(err)
	}
	err = os.MkdirAll(oldRoot, 0755)
	if err != nil {
		log.Panic(err)
	}
	err = unix.PivotRoot(newRoot, oldRoot)
	if err != nil {
		log.Panic(err)
	}
	err = os.Chdir("/")
	if err != nil {
		log.Panic(err)
	}
	err = unix.Unmount(putOld, unix.MNT_DETACH)
	if err != nil {
		log.Panic(err)
	}
	err = os.RemoveAll(putOld)
	if err != nil {
		log.Panic(err)
	}
}

func Exec() error {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(execCmd)
	defer func() {
		e.Close()
	}()
	return rootCmd.Execute()
}
