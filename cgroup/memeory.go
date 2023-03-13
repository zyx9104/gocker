package cgroup

import (
	"os"
	"path"
	"strconv"
)

const (
	memory = cgroupRoot + "/memory"
)

func MemoryTasks(uname string, pid int) {
	dir := path.Join(memory, uname)
	must(os.MkdirAll(dir, 0755))
	must(os.WriteFile(path.Join(dir, "tasks"), []byte(strconv.Itoa(pid)), 0644))
}

func Memory(uname, limit string) {
	dir := path.Join(memory, uname)
	must(os.MkdirAll(dir, 0755))
	must(os.WriteFile(path.Join(dir, "memory.limit_in_bytes"), []byte(limit), 0644))
}
