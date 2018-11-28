package cff

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"unsafe"
)

type constantPool []interface{}

func (pool constantPool) ClassInfo(index uint16) *constantClassInfo {
	return pool[index].(*constantClassInfo)
}

func (pool constantPool) FieldRefInfo(index uint16) *constantFieldRefInfo {
	return pool[index].(*constantFieldRefInfo)
}

func (pool constantPool) MethodRefInfo(index uint16) *constantMethodRefInfo {
	return pool[index].(*constantMethodRefInfo)
}

func (pool constantPool) InterfaceMethodRefInfo(index uint16) *constantInterfaceMethodRefInfo {
	return pool[index].(*constantInterfaceMethodRefInfo)
}

func (pool constantPool) StringInfo(index uint16) *constantStringInfo {
	return pool[index].(*constantStringInfo)
}

func (pool constantPool) IntegerInfo(index uint16) *constantIntegerInfo {
	return pool[index].(*constantIntegerInfo)
}

func (pool constantPool) FloatInfo(index uint16) *constantFloatInfo {
	return pool[index].(*constantFloatInfo)
}

func (pool constantPool) LongInfo(index uint16) *constantLongInfo {
	return pool[index].(*constantLongInfo)
}

func (pool constantPool) DoubleInfo(index uint16) *constantDoubleInfo {
	return pool[index].(*constantDoubleInfo)
}

func (pool constantPool) NameAndTypeInfo(index uint16) *constantNameAndTypeInfo {
	return pool[index].(*constantNameAndTypeInfo)
}

func (pool constantPool) Utf8Info(index uint16) *constantUtf8Info {
	return pool[index].(*constantUtf8Info)
}

func (pool constantPool) MethodHandleInfo(index uint16) *constantMethodHandleInfo {
	return pool[index].(*constantMethodHandleInfo)
}

func (pool constantPool) MethodTypeInfo(index uint16) *constantMethodTypeInfo {
	return pool[index].(*constantMethodTypeInfo)
}

func (pool constantPool) DynamicInfo(index uint16) *constantDynamicInfo {
	return pool[index].(*constantDynamicInfo)
}

func (pool constantPool) InvokeDynamicInfo(index uint16) *constantInvokeDynamicInfo {
	return pool[index].(*constantInvokeDynamicInfo)
}

func (pool constantPool) ModuleInfo(index uint16) *constantModuleInfo {
	return pool[index].(*constantModuleInfo)
}

func (pool constantPool) PackageInfo(index uint16) *constantPackageInfo {
	return pool[index].(*constantPackageInfo)
}

func parseConstantPool(data []byte) ([]byte, interface{}, bool) {
	switch data[0] {
	case 7:
		return data[3:], &constantClassInfo{
			binary.BigEndian.Uint16(data[1:]),
		}, false

	case 9:
		return data[5:], &constantFieldRefInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 10:
		return data[5:], &constantMethodRefInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 11:
		return data[5:], &constantInterfaceMethodRefInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 8:
		return data[3:], &constantStringInfo{
			binary.BigEndian.Uint16(data[1:]),
		}, false

	case 3:
		return data[5:], &constantIntegerInfo{
			data[1:5],
		}, false

	case 4:
		return data[5:], &constantFloatInfo{
			data[1:5],
		}, false

	case 5:
		return data[9:], &constantLongInfo{
			data[1:9],
		}, true

	case 6:
		return data[9:], &constantDoubleInfo{
			data[1:9],
		}, true

	case 12:
		return data[5:], &constantNameAndTypeInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 1:
		{
			info := constantUtf8Info{}
			length := binary.BigEndian.Uint16(data[1:])
			offset := 3 + int(length)
			info.data = data[3:offset]
			return data[offset:], &info, false
		}

	case 15:
		return data[4:], &constantMethodHandleInfo{
			data[1],
			binary.BigEndian.Uint16(data[2:]),
		}, false

	case 16:
		return data[3:], &constantMethodTypeInfo{
			binary.BigEndian.Uint16(data[1:]),
		}, false

	case 17:
		return data[5:], &constantDynamicInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 18:
		return data[5:], &constantInvokeDynamicInfo{
			binary.BigEndian.Uint16(data[1:]),
			binary.BigEndian.Uint16(data[3:]),
		}, false

	case 19:
		return data[3:], &constantModuleInfo{
			binary.BigEndian.Uint16(data[1:]),
		}, false

	case 20:
		return data[3:], &constantPackageInfo{
			binary.BigEndian.Uint16(data[1:]),
		}, false
	}

	panic(fmt.Errorf("unknown constant type: %d", data[0]))
}

type constantClassInfo struct {
	nameIndex uint16
}

type constantFieldRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantMethodRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantInterfaceMethodRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantStringInfo struct {
	stringIndex uint16
}

type constantIntegerInfo struct {
	data []byte // 长度 4
}

func (info *constantIntegerInfo) Int32() int32 {
	return *((*int32)(unsafe.Pointer(&info.data[0])))
}

type constantFloatInfo struct {
	data []byte // 长度 4
}

func (info *constantFloatInfo) Float32() float32 {
	return *((*float32)(unsafe.Pointer(&info.data[0])))
}

type constantLongInfo struct {
	data []byte // 长度 8
}

func (info *constantLongInfo) Int64() int64 {
	return *((*int64)(unsafe.Pointer(&info.data[0])))
}

type constantDoubleInfo struct {
	data []byte // 长度 8
}

func (info *constantDoubleInfo) Float64() float64 {
	return *((*float64)(unsafe.Pointer(&info.data[0])))
}

type constantNameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

type constantUtf8Info struct {
	data []byte // 长度 length
}

func (utf8 *constantUtf8Info) String() string {
	size := len(utf8.data)
	index := 0
	buf := make([]uint16, 0, 256)
	var c, cc uint16
	st := 0

	for index < size {
		c = uint16(utf8.data[index])
		index++
		switch st {
		case 0:
			if c < 0x80 {
				buf = append(buf, c)
			} else if c < 0xE0 && c > 0xBF {
				cc = c & 0x1F
				st = 1
			} else {
				cc = c & 0x0F
				st = 2
			}

		case 1:
			buf = append(buf, (cc<<6)|(c&0x3F))
			st = 0

		case 2:
			cc = (cc << 6) | (c & 0x3F)
			st = 1
		}
	}

	return string(utf16.Decode(buf))
}

type constantMethodHandleInfo struct {
	referenceKind  byte
	referenceIndex uint16
}

type constantMethodTypeInfo struct {
	descriptorIndex uint16
}

type constantDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

type constantInvokeDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

type constantModuleInfo struct {
	nameIndex uint16
}

type constantPackageInfo struct {
	nameIndex uint16
}
