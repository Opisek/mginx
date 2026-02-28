package main

import (
	"mginx/config"
	"mginx/connections/server"
	"mginx/models"
	"mginx/protocol/payloads"
	"mginx/util"
)

func main() {
	conf := config.ReadConfig()

	packetQueue := make(chan util.Pair[*models.GameClient, payloads.GenericPacket])

	go server.HandlePackets(packetQueue, conf)
	server.StartServer("localhost", 25565, packetQueue, conf)
}
