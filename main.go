package main

import (
	"path"
	"fmt"
	"lsm-verification/config"
)

func main() {
	conf := config.LoadConfig(path.Join("config", "config.yaml"))
	fmt.Println(conf)
}
