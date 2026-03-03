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

	// Watch managed servers
	for _, server := range conf.Servers {
		if server.Watchdog.IsManaged() {
			go watchdog.WatchUpstream(server)
		}
	}
	time.Sleep(3 * time.Second) // Let the watchdog initialize for every server

	// Channel for fully buffered packets to be processed further
	packetQueue := make(chan util.Pair[*models.DownstreamClient, payloads.GenericPacket])

	// Handle fully buffered packets
	go downstream.HandlePackets(packetQueue, conf)

	// Handle connections
	downstream.StartServer("localhost", 25565, packetQueue, conf)
}
