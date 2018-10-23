package main

import (
	"IT-Berries_Go_server/gameServer"
)

func main() {
	var server gameServer.TheServer
	//Server init
	server.Prepare()
	server.Start()
}
