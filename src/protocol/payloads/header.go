package payloads

type GenericPacket struct {
	Length       uint64
	Id           uint64
	Payload      []byte
	ActualLength uint64
}
