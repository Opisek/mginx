package server

import (
	"bytes"
	"errors"
	"fmt"
	"mginx/config"
	"mginx/models"
	"mginx/protocol"
	"mginx/protocol/parsing"
	"mginx/protocol/payloads"
	"mginx/util"
	"net"
	"time"
)

func StartServer(address string, port uint16, packetQueue chan util.Pair[*models.GameClient, payloads.GenericPacket], conf *config.Configuration) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", address, port))
	if err != nil {
		return errors.Join(fmt.Errorf("could not start listening on %v:%v", address, port), err)
	}
	defer listener.Close()

	fmt.Printf("mginx listening on %v:%v\n", address, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClientConnection(conn, packetQueue, conf)
	}
}

func handleClientConnection(conn net.Conn, packetQueue chan util.Pair[*models.GameClient, payloads.GenericPacket], conf *config.Configuration) {
	defer conn.Close()

	client := &models.GameClient{
		Connection: conn,
	}

	var buffer bytes.Buffer

	data := make([]byte, 1024)
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(data)
		if !client.IsAlive() {
			return
		}
		if err != nil {
			client.Kill()
			return
		}
		if client.IsProxying() {
			client.UpstreamConnection.Write(data[:n])
			continue
		}

		buffer.Write(data[:n])

		packet, err := parsing.ParseHeader(buffer.Bytes())
		if err != nil {
			continue
		}
		if packet.Length > packet.ActualLength {
			continue
		}

		remainingPayload := packet.Payload
		cutIndex := len(remainingPayload) - int(packet.ActualLength-packet.Length)

		packet.Payload = make([]byte, cutIndex)
		copy(packet.Payload, remainingPayload[:cutIndex])

		buffer.Reset()
		buffer.Write(remainingPayload[cutIndex:])

		packetQueue <- util.Pair[*models.GameClient, payloads.GenericPacket]{
			First:  client,
			Second: packet,
		}
	}
}

func HandlePackets(packetQueue chan util.Pair[*models.GameClient, payloads.GenericPacket], conf *config.Configuration) {
	for {
		received := <-packetQueue
		client := received.First
		packet := received.Second

		if !client.IsAlive() || client.IsProxying() {
			continue
		}

		err := protocol.HandlePacket(client, packet, conf)
		if err != nil {
			fmt.Println(errors.Join(errors.New("could not handle client packet"), err))

			client.Kill()
		}
	}
}
