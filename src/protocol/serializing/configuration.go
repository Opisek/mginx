package serializing

import (
	"bytes"
	util "mginx/protocol/internal"
	"mginx/protocol/payloads"
)

func SerializeTransfer(payload payloads.Transfer) []byte {
	var buffer bytes.Buffer

	buffer.Write(util.SerializeString(payload.Host))
	buffer.Write(util.SerializeVarInt(uint64(payload.Port)))

	return SerializePacketWithHeader(0x0B, buffer.Bytes())
}
