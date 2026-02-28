package main

import (
	"mginx/config"
	"mginx/connections/downstream"
	"mginx/connections/watchdog"
	"mginx/models"
	"mginx/protocol/payloads"
	"mginx/util"
)

func main() {
	conf := config.ReadConfig()

	for _, server := range conf.Servers {
		go watchdog.WatchUpstream(server)
	}

	packetQueue := make(chan util.Pair[*models.DownstreamClient, payloads.GenericPacket])

	go downstream.HandlePackets(packetQueue, conf)
	downstream.StartServer("localhost", 25565, packetQueue, conf)
}
