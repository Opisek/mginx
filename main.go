package main

import (
	"bytes"
	"errors"
	"fmt"
	"mginx/models"
	"mginx/protocol"
	"mginx/protocol/parsing"
	"mginx/util"
	"net"
	"time"
)

func handleConnection(conn net.Conn, packetQueue chan util.Pair[*models.GameClient, parsing.GenericPacket]) {
	defer conn.Close()

	client := &models.GameClient{
		Connection: conn,
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var buffer bytes.Buffer

	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
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

		packetQueue <- util.Pair[*models.GameClient, parsing.GenericPacket]{
			First:  client,
			Second: packet,
		}
	}
}

func handleCompletePackets(packetQueue chan util.Pair[*models.GameClient, parsing.GenericPacket]) {
	for {
		received := <-packetQueue
		client := received.First
		packet := received.Second

		if client.Connection == nil {
			continue
		}

		err := protocol.HandlePacket(client, packet)
		if err != nil {
			fmt.Println(errors.Join(errors.New("could not handle client packet"), err))
			client.Connection.Close()
			client.Connection = nil
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:25565")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	packetQueue := make(chan util.Pair[*models.GameClient, parsing.GenericPacket])

	go handleCompletePackets(packetQueue)

	fmt.Println("Server running on :25565")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn, packetQueue)
	}
}
