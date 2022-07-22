package main

import (
	"runtime"
	"x-game/client/_client"
	"x-game/client/config"
	"x-game/client/handlers"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handlers.InitHandlers()
	_client.Client = _client.NewClient(config.GetClientConfig().ServerUrl)
	_client.Client.Run()
}
