package interp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ClarkGuan/jnif/cff"
)

type JavaType string

func (jt JavaType) JavaType() (s string) {
	index := 0

	for index < len(jt) {
		switch jt[index] {
		case 'V':
			return "void"

		case 'B':
			return "byte" + s

		case 'C':
			return "char" + s

		case 'D':
			return "double" + s

		case 'F':
			return "float" + s

		case 'I':
			return "int" + s

		case 'J':
			return "long" + s

		case 'S':
			return "short" + s

		case 'Z':
			return "boolean" + s

		case '[':
			s += "[]"
			index++

		default:
			return fmt.Sprintf("%s%s",
				strings.Replace(string(jt[index+1:len(jt)-1]), "/", ".", -1), s)
		}
	}

	panic(fmt.Errorf("can't reach here"))
}

func (jt JavaType) JavaZero() string {
	switch jt[0] {
	case 'V':
		return ""

	case 'Z':
		return "false"

	case '[':
		return "null"

	case 'L':
		return "null"

	default:
		return "0"
	}

	panic(fmt.Errorf("can't reach here"))
}

func (jt JavaType) IsVoid() bool {
	return jt[0] == 'V'
}

func (jt JavaType) IsArray() bool {
	return jt[0] == '['
}

func (jt JavaType) IsObject() bool {
	return jt[0] == 'L'
}

func (jt JavaType) IsNormal() bool {
	return !jt.IsArray() && !jt.IsObject() && !jt.IsVoid()
}

type Method struct {
	Name       string
	ClassName  string
	Arguments  []JavaType
	ReturnType JavaType
	Desc       string
	cff.Modifier
	overload int // 已有同名函数的个数。0 表示前面没有任何同名函数
}

func (m *Method) JavaClassName() string {
	return strings.Replace(m.ClassName, "/", ".", -1)
}

func (m *Method) JavaSignature() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s ", m.Modifier)
	fmt.Fprintf(buf, "%s %s(", m.ReturnType.JavaType(), m.Name)

	for n, r := range m.Arguments {
		if n > 0 {
			fmt.Fprintf(buf, ", ")
		}
		fmt.Fprintf(buf, "%s arg%d", r.JavaType(), n+1)
	}

	fmt.Fprintf(buf, ")")
	return buf.String()
}

func (m *Method) OverloadName() string {
	if m.overload == 0 {
		return m.Name
	}
	return fmt.Sprintf("%s%d", m.Name, m.overload)
}

var noRegisterErr = errors.New("no register")

var gRegisters = make(map[string]Transformer)

func Register(name string, trans Transformer) {
	gRegisters[name] = trans
}

type Transformer interface {
	Init(args []string) (string, error)

	Transform(methods map[string][]*Method, maxCount int) error
}
