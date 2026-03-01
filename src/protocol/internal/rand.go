package util

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
)

func GetRandomLong() (uint64, error) {
	var data [8]byte
	_, err := rand.Read(data[:])
	if err != nil {
		return 0, errors.Join(errors.New("coud not generate a random number"), err)
	}
	return binary.LittleEndian.Uint64(data[:]), nil
}
