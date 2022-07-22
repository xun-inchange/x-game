package main

import (
	"runtime"
	"x-game/server/config"
	"x-game/server/modles/server"
	"x-game/server/module/handler"
	network "x-game/x-common/x_net"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handler.InitHandlers()
	s := network.NewServer(config.GetServerConfig().App.Addr)
	server.SeverInfo.SetNetServer(s)
	server.SeverInfo.Run()
}
