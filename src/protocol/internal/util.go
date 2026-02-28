package util

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

func ParseVarInt(buffer []byte) (uint64, []byte, error) {
	len, read := binary.Uvarint(buffer)

	if read == 0 {
		return 0, nil, errors.New("improperly formatted varint")
	}

	return len, buffer[read:], nil
}

func SerializeVarInt(num uint64) []byte {
	buffer := make([]byte, 128)
	n := binary.PutUvarint(buffer, num)

	return buffer[:n]
}

func ParseString(buffer []byte) (string, []byte, error) {
	size, buffer, err := ParseVarInt(buffer)

	if err != nil {
		return "", nil, errors.Join(errors.New("could not parse string length"), err)
	}

	if uint64(len(buffer)) < size {
		return "", nil, errors.New("remaining buffer too short to store string of specified size")
	}

	return string(buffer[:size]), buffer[size:], nil
}

func SerializeString(str string) []byte {
	var buffer bytes.Buffer

	buffer.Write(SerializeVarInt(uint64(len(str))))
	buffer.WriteString(str)

	return buffer.Bytes()
}

func ParseUnsignedShort(buffer []byte) (uint16, []byte, error) {
	if len(buffer) < 2 {
		return 0, nil, errors.New("remaining buffer too short to store an unsigned short")
	}

	return binary.BigEndian.Uint16(buffer[:2]), buffer[2:], nil
}

func SerializeUnsignedShort(num uint16) []byte {
	return binary.BigEndian.AppendUint16(nil, num)
}

func ParseUuid(buffer []byte) (uuid.UUID, []byte, error) {
	if len(buffer) < 16 {
		return uuid.Nil, nil, errors.New("remaining buffer too short to store a uuid")
	}

	id, err := uuid.FromBytes(buffer[:16])

	if err != nil {
		return uuid.Nil, nil, errors.Join(errors.New("could not parse uuid"), err)
	}

	return id, buffer[16:], nil
}

func SerializeUuid(id uuid.UUID) []byte {
	bytes, err := id.MarshalBinary()
	if err != nil {
		panic(errors.Join(errors.New("could not serialize uuid"), err))
	}

	return bytes
}
