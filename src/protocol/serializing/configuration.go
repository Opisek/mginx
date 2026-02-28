package serializing

import (
	"bytes"
	util "mginx/protocol/internal"
)

type TransferPayload struct {
	Host string
	Port uint16
}

func SerializeTransfer(payload TransferPayload) []byte {
	var buffer bytes.Buffer

	buffer.Write(util.SerializeString(payload.Host))
	buffer.Write(util.SerializeVarInt(uint64(payload.Port)))

	return SerializePacketWithHeader(0x0B, buffer.Bytes())
}
