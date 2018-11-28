package cff

import (
	"encoding/binary"
)

type fieldInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []*attributeInfo // 长度 attributesCount
}

func parseFieldInfo(data []byte) ([]byte, *fieldInfo) {
	fi := fieldInfo{}

	fi.accessFlags = binary.BigEndian.Uint16(data)
	data = data[2:]

	fi.nameIndex = binary.BigEndian.Uint16(data)
	data = data[2:]

	fi.descriptorIndex = binary.BigEndian.Uint16(data)
	data = data[2:]

	attributesCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var ai *attributeInfo
	for i := 0; i < int(attributesCount); i++ {
		data, ai = parseAttributeInfo(data)
		fi.attributes = append(fi.attributes, ai)
	}

	return data, &fi
}

type methodInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []*attributeInfo // 长度 attributesCount
}

func parseMethodInfo(data []byte) ([]byte, *methodInfo) {
	mi := methodInfo{}

	mi.accessFlags = binary.BigEndian.Uint16(data)
	data = data[2:]

	mi.nameIndex = binary.BigEndian.Uint16(data)
	data = data[2:]

	mi.descriptorIndex = binary.BigEndian.Uint16(data)
	data = data[2:]

	attributesCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var ai *attributeInfo
	for i := 0; i < int(attributesCount); i++ {
		data, ai = parseAttributeInfo(data)
		mi.attributes = append(mi.attributes, ai)
	}

	return data, &mi
}

type classFile struct {
	magic        uint32
	minorVersion uint16
	majorVersion uint16
	constantPool // 长度 constantPoolCount - 1
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16         // 长度 interfacesCount
	fields       []*fieldInfo     // 长度 fieldsCount
	methods      []*methodInfo    // 长度 methodsCount
	attributes   []*attributeInfo // 长度 attributesCount
}

func parseClassFile(data []byte) (cf *classFile, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	cf = new(classFile)

	cf.magic = binary.BigEndian.Uint32(data)
	data = data[4:]

	cf.minorVersion = binary.BigEndian.Uint16(data)
	data = data[2:]

	cf.majorVersion = binary.BigEndian.Uint16(data)
	data = data[2:]

	constantPoolCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var cons interface{}
	cf.constantPool = make(constantPool, constantPoolCount)
	var is64 bool
	for i := 1; i < int(constantPoolCount); i++ {
		data, cons, is64 = parseConstantPool(data)
		cf.constantPool[i] = cons
		if is64 { // long 和 double 占用两个索引号
			i++
		}
	}

	cf.accessFlags = binary.BigEndian.Uint16(data)
	data = data[2:]

	cf.thisClass = binary.BigEndian.Uint16(data)
	data = data[2:]

	cf.superClass = binary.BigEndian.Uint16(data)
	data = data[2:]

	interfacesCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	for i := 0; i < int(interfacesCount); i++ {
		cf.interfaces = append(cf.interfaces, binary.BigEndian.Uint16(data))
		data = data[2:]
	}

	fieldsCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var fi *fieldInfo
	for i := 0; i < int(fieldsCount); i++ {
		data, fi = parseFieldInfo(data)
		cf.fields = append(cf.fields, fi)
	}

	methodsCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var mi *methodInfo
	for i := 0; i < int(methodsCount); i++ {
		data, mi = parseMethodInfo(data)
		cf.methods = append(cf.methods, mi)
	}

	attributesCount := binary.BigEndian.Uint16(data)
	data = data[2:]

	var ai *attributeInfo
	for i := 0; i < int(attributesCount); i++ {
		data, ai = parseAttributeInfo(data)
		cf.attributes = append(cf.attributes, ai)
	}

	return
}
