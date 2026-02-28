package payloads

type Handshake struct {
	Version uint64
	Address string
	Port    uint16
	Intent  uint64
}
