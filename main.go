package main

import (
	"log"
	"lsm-verification/config"
	"path"
)

func main() {
	config := config.LoadConfig(path.Join("config", "config.yaml"))
	log.Println("Running in mode: ", config.RunMode)

}
