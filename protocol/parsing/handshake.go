package parsing

import (
	"errors"
	"fmt"
	util "mginx/protocol/internal"
)

type HandshakePayload struct {
	Version uint64
	Address string
	Port    uint16
	Intent  uint64
}

func ParseHandshake(buffer []byte) (HandshakePayload, error) {
	version, buffer, err := util.ParseVarInt(buffer)
	if err != nil {
		return HandshakePayload{},
			errors.Join(errors.New("could not parse client version"), err)
	}

	address, buffer, err := util.ParseString(buffer)
	if err != nil {
		return HandshakePayload{},
			errors.Join(errors.New("could not parse client address"), err)
	}

	port, buffer, err := util.ParseUnsignedShort(buffer)
	if err != nil {
		return HandshakePayload{},
			errors.Join(errors.New("could not parse client port"), err)
	}

	intent, buffer, err := util.ParseVarInt(buffer)
	if err != nil {
		return HandshakePayload{},
			errors.Join(errors.New("could not parse client intent"), err)
	}

	if len(buffer) != 0 {
		return HandshakePayload{},
			fmt.Errorf("remaining bytes at the end of payload: %v", len(buffer))
	}

	return HandshakePayload{
		Version: version,
		Address: address,
		Port:    port,
		Intent:  intent,
	}, nil
}
