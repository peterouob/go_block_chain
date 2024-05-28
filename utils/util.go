package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func ToHexint(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Println("Error to translate num to []byte :", err)
		return nil
	}
	return buff.Bytes()
}
