package main

import (
	"github.com/lfroomin/restaurant-container/config"
	"github.com/lfroomin/restaurant-container/server"
)

func main() {
	appCfg := config.Init(".")
	server.Init(appCfg)
}
