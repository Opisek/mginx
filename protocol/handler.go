package protocol

import (
	"errors"
	"fmt"
	"mginx/models"
	"mginx/protocol/parsing"
	"mginx/protocol/phases"
)

func HandlePacket(client *models.GameClient, packet parsing.GenericPacket) error {
	switch client.GamePhase {
	case 0x00:
		err := phases.HandleHandshakePhase(client, packet)
		if err != nil {
			return errors.Join(errors.New("could not handle packet in handshake phase"), err)
		}
	case 0x01:
		fmt.Println("status phase")
	case 0x02:
		err := phases.HandleLoginPhase(client, packet)
		if err != nil {
			return errors.Join(errors.New("could not handle packet in login phase"), err)
		}
	case 0x04:
		fmt.Println("configuration phase")
	default:
		return fmt.Errorf("invalid game phase: %v", client.GamePhase)
	}

	return nil
}
