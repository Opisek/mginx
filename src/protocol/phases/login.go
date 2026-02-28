package phases

import (
	"errors"
	"fmt"
	"mginx/config"
	"mginx/connections/upstream"
	"mginx/models"
	"mginx/protocol/parsing"
	"mginx/protocol/payloads"
	"mginx/protocol/serializing"
)

func HandleLoginPhase(client *models.DownstreamClient, packet payloads.GenericPacket, conf *config.Configuration) error {
	switch packet.Id {
	case 0x00:
		err := handleClientLoginStart(client, packet)
		if err != nil {
			return errors.Join(errors.New("could not parse login start packet"), err)
		}
	case 0x03:
		err := handleClientLoginAcknowledged(client, packet)
		if err != nil {
			return errors.Join(errors.New("could not parse login acknowledged packet"), err)
		}
	default:
		return fmt.Errorf("invalid packet id: %v", packet.Id)
	}
	return nil
}

func handleClientLoginStart(client *models.DownstreamClient, packet payloads.GenericPacket) error {
	payload, err := parsing.ParseLoginStart(packet.Payload)

	if err != nil {
		return err
	}

	client.Username = payload.Name
	client.Uuid = payload.Uuid

	// If we are working with redirects, acknowledge to move into the configuration phase
	if client.Upstream.Redirect {
		client.Connection.Write(serializing.SerializeLoginSuccess(payloads.LoginSuccess{
			Name: client.Username,
			Uuid: client.Uuid,
		}))
		return nil
	}

	// Otherwise, set up a proxy channel
	actualAddress, err := upstream.ProxyConnection(client)
	if err != nil {
		return errors.Join(errors.New("could not proxy connection"), err)
	}

	data := serializing.SerializeHandshake(payloads.Handshake{
		Version: client.Version,
		Address: actualAddress,
		Port:    client.Port,
		Intent:  0x02,
	})
	client.UpstreamConnection.Write(data)

	data = serializing.SerializeLoginStart(payload)
	client.UpstreamConnection.Write(data)

	return nil
}

func handleClientLoginAcknowledged(client *models.DownstreamClient, packet payloads.GenericPacket) error {
	_, err := parsing.ParseLoginAcknowledged(packet.Payload)

	if err != nil {
		return err
	}

	client.GamePhase = 0x04

	client.Connection.Write(serializing.SerializeTransfer(payloads.Transfer{
		Host: client.Upstream.To.Hostname,
		Port: client.Upstream.To.Port,
	}))
	client.StartTransfer()

	return nil
}
