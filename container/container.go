package container

import (
	"log"
	"path"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

const (
	BasePath       = "/var/lib/gocker"
	BaseImage      = BasePath + "/images"
	BaseContainers = BasePath + "/containers"
	BaseDiff       = BasePath + "/diff"
	Config         = BasePath + "/config.json"
)

type Container struct {
	Id          string `json:"id"`
	Pid         int    `json:"pid"`
	Name        string `json:"name"`
	WriterLayer string `json:"writer_layer"`
	BaseImage   string `json:"base_image"`
	Volume      Volume `json:"volume"`
	Status      string `json:"status"`
}

type Volume struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func New(BaseImage, name string) *Container {
	write := uuid.NewString()
	id := uuid.NewString()
	if name == "" {
		name = id
	}
	return &Container{
		Id:          id,
		Name:        name,
		WriterLayer: write,
		BaseImage:   BaseImage,
	}
}

func must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (c *Container) Close() {
	must(unix.Unmount(path.Join(BaseContainers, c.Volume.Target), 0))
	must(unix.Unmount(path.Join(BaseContainers, c.Id), 0))
}
