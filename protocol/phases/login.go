package phases

import (
	"errors"
	"fmt"
	"mginx/models"
	"mginx/protocol/parsing"
	"mginx/protocol/serializing"
)

func HandleLoginPhase(client *models.GameClient, packet parsing.GenericPacket) error {
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

func handleClientLoginStart(client *models.GameClient, packet parsing.GenericPacket) error {
	payload, err := parsing.ParseLoginStart(packet.Payload)

	if err != nil {
		return err
	}

	fmt.Printf("login start: %v, %v\n", payload.Name, payload.Uuid)

	client.Username = payload.Name
	client.Uuid = payload.Uuid

	client.Connection.Write(serializing.SerializeLoginSuccess(serializing.LoginSuccessPayload{
		Name: client.Username,
		Uuid: client.Uuid,
	}))

	return nil
}

func handleClientLoginAcknowledged(client *models.GameClient, packet parsing.GenericPacket) error {
	_, err := parsing.ParseLoginAcknowledged(packet.Payload)

	if err != nil {
		return err
	}

	client.GamePhase = 0x04

	client.Connection.Write(serializing.SerializeTransfer(serializing.TransferPayload{
		Host: "example.com",
		Port: 25565,
	}))

	fmt.Printf("login acknowledged\n")

	return nil
}
