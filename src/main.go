package main

import (
	"mginx/config"
	"mginx/connections/downstream"
	"mginx/connections/watchdog"
	"mginx/models"
	"mginx/protocol/payloads"
	"mginx/util"
	"time"
)

func main() {
	conf := config.ReadConfig()

	for _, server := range conf.Servers {
		if server.Watchdog.IsManaged() {
			go watchdog.WatchUpstream(server)
		}
	}
	time.Sleep(2 * time.Second) // Let the watchdog initialize for every server

	packetQueue := make(chan util.Pair[*models.DownstreamClient, payloads.GenericPacket])

	go downstream.HandlePackets(packetQueue, conf)
	downstream.StartServer("localhost", 25565, packetQueue, conf)
}
