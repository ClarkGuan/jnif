package cff

import (
	"encoding/binary"
)

type attributeInfo struct {
	nameIndex uint16
	info      []byte
}

func parseAttributeInfo(data []byte) ([]byte, *attributeInfo) {
	ai := attributeInfo{}

	ai.nameIndex = binary.BigEndian.Uint16(data)
	data = data[2:]

	length := binary.BigEndian.Uint32(data)
	data = data[4:]

	ai.info = data[:length]
	data = data[length:]

	return data, &ai
}
