package phases

import (
	"errors"
	"fmt"
	"mginx/config"
	"mginx/models"
	util "mginx/protocol/internal"
	"mginx/protocol/parsing"
	"mginx/protocol/payloads"
	"mginx/protocol/serializing"
	"time"
)

func HandleConfigurationPhase(client *models.DownstreamClient, packet payloads.GenericPacket, conf *config.Configuration) error {
	switch packet.Id {
	case 0x04:
		err := handleClientKeepAlive(client, packet)
		if err != nil {
			return errors.Join(errors.New("could not parse keepalive packet"), err)
		}
	default:
		fmt.Printf("unknown config packet id: %v", packet.Id)
		return nil
		//return fmt.Errorf("invalid packet id: %v", packet.Id)
	}
	return nil
}

func handleClientKeepAlive(client *models.DownstreamClient, packet payloads.GenericPacket) error {
	payload, err := parsing.ParseKeepAlive(packet.Payload)

	if err != nil {
		return err
	}

	if payload.Id != client.ExpectedKeepalive {
		return fmt.Errorf("keepalive id %v received but expected %v", payload.Id, client.ExpectedKeepalive)
	}

	fmt.Println("received keepalive", payload.Id)

	client.ExpectedKeepalive, err = util.GetRandomLong()
	if err != nil {
		return errors.Join(errors.New("could not generate a random keepalive ID"), err)
	}

	go func() {
		time.Sleep(5 * time.Second)

		if !client.IsAlive() || client.IsProxying() {
			return
		}

		client.Connection.Write(serializing.SerializeKeepAlive(payloads.KeepAlive{
			Id: client.ExpectedKeepalive,
		}))
	}()

	return nil
}
