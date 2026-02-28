package phases

import (
	"errors"
	"fmt"
	"mginx/config"
	"mginx/models"
	"mginx/protocol/parsing"
)

func HandleHandshakePhase(client *models.GameClient, packet parsing.GenericPacket, conf *config.Configuration) error {
	switch packet.Id {
	case 0x00:
		err := handleClientHandshake(client, packet, conf)
		if err != nil {
			return errors.Join(errors.New("could not parse handshake packet"), err)
		}
	default:
		return fmt.Errorf("invalid packet id: %v", packet.Id)
	}
	return nil
}

func handleClientHandshake(client *models.GameClient, packet parsing.GenericPacket, conf *config.Configuration) error {
	payload, err := parsing.ParseHandshake(packet.Payload)

	if err != nil {
		return err
	}

	fmt.Printf("handshake received: %v, %v, %v, %v\n", payload.Version, payload.Address, payload.Port, payload.Intent)

	switch payload.Intent {
	case 0x01:
		client.GamePhase = 0x01
	case 0x02:
		fallthrough
	case 0x03:
		client.GamePhase = 0x02
	default:
		return fmt.Errorf("invalid handshake intent: %v", payload.Intent)
	}

	client.Version = payload.Version
	client.Address = payload.Address
	client.Port = payload.Port
	client.Upstream = conf.GetUpstream(client.Address, client.Port)

	if client.Upstream == nil {
		return fmt.Errorf(`no upstream found for address "%v:%v"`, client.Address, client.Port)
	}

	return nil
}
