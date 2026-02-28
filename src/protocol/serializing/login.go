package serializing

import (
	"bytes"
	util "mginx/protocol/internal"

	"github.com/google/uuid"
)

type LoginSuccessPayload struct {
	Name string
	Uuid uuid.UUID
}

func SerializeLoginSuccess(payload LoginSuccessPayload) []byte {
	var buffer bytes.Buffer

	buffer.Write(util.SerializeUuid(payload.Uuid))
	buffer.Write(util.SerializeString(payload.Name))
	buffer.Write(util.SerializeVarInt(0))

	return SerializePacketWithHeader(0x02, buffer.Bytes())
}
