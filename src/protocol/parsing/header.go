package parsing

import (
	util "mginx/protocol/internal"
)

type GenericPacket struct {
	Length       uint64
	Id           uint64
	Payload      []byte
	ActualLength uint64
}

func ParseHeader(buffer []byte) (GenericPacket, error) {
	packetLen, buffer, err := util.ParseVarInt(buffer)
	if err != nil {
		return GenericPacket{}, err
	}

	actualLen := uint64(len(buffer))

	packetId, buffer, err := util.ParseVarInt(buffer)
	if err != nil {
		return GenericPacket{}, err
	}

	return GenericPacket{
		Length:       packetLen,
		Id:           packetId,
		Payload:      buffer,
		ActualLength: actualLen,
	}, nil
}
