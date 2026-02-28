package watchdog

import (
	"bytes"
	"errors"
	"fmt"
	"mginx/connections/upstream"
	"mginx/constants"
	"mginx/models"
	"mginx/protocol/parsing"
	"mginx/protocol/payloads"
	"mginx/protocol/serializing"
	"net"
	"time"
)

func handleUpstreamStatusConnection(conn net.Conn, res chan int) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var buffer bytes.Buffer

	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			res <- -1
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

		statusResponse, err := parsing.ParseStatusResponse(packet.Payload[:len(packet.Payload)-int(packet.ActualLength-packet.Length)])
		if err != nil {
			res <- -1
			return
		}

		res <- statusResponse.Players.Online
		return
	}
}

func checkStatus(server *models.UpstreamServer) (int, error) {
	res := make(chan int)

	conn, address, err := upstream.StartClient(server.To.Hostname, server.To.Port, func(conn net.Conn) {
		handleUpstreamStatusConnection(conn, res)
	})

	if err != nil {
		return -1, errors.Join(errors.New("could not check server status"), err)
	}

	conn.Write(serializing.SerializeHandshake(payloads.Handshake{
		Version: constants.ProtocolVersion,
		Address: address,
		Port:    server.To.Port,
		Intent:  0x01,
	}))

	conn.Write(serializing.SerializeStatusRequest(payloads.StatusRequest{}))

	return <-res, nil
}

func WatchUpstream(server *models.UpstreamServer) {
	todo := make(chan bool)

	for {
		select {
		case <-time.After(5 * time.Second):
		case <-todo:
		}

		players, err := checkStatus(server)

		if err != nil || players == -1 {
			players = 0
		}

		fmt.Printf("%s has %v online players\n", server.InternalName, players)
	}
}
