package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"gocker/container"
)

type Engine struct {
	Containers map[string]*container.Container
}

func New() *Engine {

	bs, err := os.ReadFile(container.Config)
	if err != nil {
		log.Panic(err)
	}
	engine := &Engine{
		Containers: make(map[string]*container.Container),
	}
	if err = json.Unmarshal(bs, &engine.Containers); err != nil {
		log.Panic(err)
	}
	log.Infof("load %d containers", len(engine.Containers))
	return engine
}

func (e *Engine) CreateContainer(image string) *container.Container {
	c := container.New(image, "")
	e.Containers[c.Id] = c
	return c
}

func (e *Engine) RunContainer(image string) *container.Container {
	c := e.CreateContainer(image)
	writerPath := path.Join(container.BaseDiff, c.WriterLayer)
	containerPath := path.Join(container.BaseContainers, c.Id)
	err := os.MkdirAll(containerPath, 0755)
	if err != nil {
		log.Panic(err)
	}
	err = os.MkdirAll(writerPath, 0755)
	if err != nil {
		log.Panic(err)
	}
	err = e.MountAufs(containerPath, writerPath, []string{path.Join(container.BaseImage, c.BaseImage)})
	if err != nil {
		log.Panic(err)
	}
	return c
}

func (e *Engine) MountAufs(target, rwLayers string, layers []string) error {
	roBranch := ""
	for _, s := range layers {
		roBranch += fmt.Sprintf("%v=ro:", s)
	}
	rw := fmt.Sprintf("%v=rw", rwLayers)
	branches := fmt.Sprintf("br:%v:%v", rw, roBranch)

	log.Infof("mount aufs target %v, branches: %v", target, branches)
	return unix.Mount("none", target, "aufs", 0, branches)
}

func (e *Engine) Close() {
	bs, err := json.Marshal(e.Containers)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(container.Config, bs, 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (e *Engine) GetContainer(key string) *container.Container {
	for k, v := range e.Containers {
		if strings.HasPrefix(k, key) || v.Name == key {
			return v
		}
	}
	return nil
}
