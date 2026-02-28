package serializing

import "mginx/protocol/payloads"

func SerializeStatusRequest(payloads payloads.StatusRequest) []byte {
	return SerializePacketWithHeader(0x00, []byte{})
}
