package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func ToHexint(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Println("Error to translate num to []byte :", err)
		return nil
	}
	return buff.Bytes()
}

func FileExsit(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

func Handle(err error) {
	if err != nil {
		log.Println("Error :", err)
	}
}
