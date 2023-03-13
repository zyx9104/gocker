//go:build linux

package main

import (
	"gocker/cmd"
	"gocker/log"
	_ "gocker/setns"
)

func main() {
	log.Info("start")

	if err := cmd.Exec(); err != nil {
		log.Panic(err)
	}
}
