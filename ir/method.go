package ir

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ClarkGuan/jnif/cff"
)

type javaType string

func (jt javaType) JavaType() (s string) {
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

func (jt javaType) CType() string {
	switch jt[0] {
	case 'V':
		return "void"

	case 'B':
		return "GoUint8"

	case 'C':
		return "GoUint16"

	case 'D':
		return "GoFloat64"

	case 'F':
		return "GoFloat32"

	case 'I':
		return "GoInt32"

	case 'J':
		return "GoInt64"

	case 'S':
		return "GoInt16"

	case 'Z':
		return "GoUint8"

	default:
		return "GoUintptr"
	}

	panic(fmt.Errorf("can't reach here"))
}

func (jt javaType) GolangType() string {
	switch jt[0] {
	case 'V':
		return ""

	case 'B':
		return "uint8"

	case 'C':
		return "uint16"

	case 'D':
		return "float64"

	case 'F':
		return "float32"

	case 'I':
		return "int32"

	case 'J':
		return "int64"

	case 'S':
		return "int16"

	case 'Z':
		return "uint8"

	default:
		return "uintptr"
	}

	panic(fmt.Errorf("can't reach here"))
}

func (jt javaType) IsVoid() bool {
	return jt[0] == 'V'
}

func (jt javaType) IsArray() bool {
	return jt[0] == '['
}

func (jt javaType) IsObject() bool {
	return jt[0] == 'L'
}

func (jt javaType) IsNormal() bool {
	return !jt.IsArray() && !jt.IsObject() && !jt.IsVoid()
}

func (jt javaType) Zero() string {
	switch jt[0] {
	case 'V':
		return ""

	case 'Z':
		return "false"

	default:
		return "0"
	}

	panic(fmt.Errorf("can't reach here"))
}

type Method struct {
	Name       string
	className  string
	arguments  []javaType
	returnType javaType
	Desc       string
	cff.Modifier
	overload int // 已有同名函数的个数。0 表示前面没有任何同名函数
}

func (m *Method) JavaClassName() string {
	return strings.Replace(m.className, "/", ".", -1)
}

func (m *Method) GoFuncName() string {
	name := m.Name
	if m.overload > 0 {
		name = fmt.Sprintf("%s%d", name, m.overload)
	}
	return fmt.Sprintf("%s_%s",
		strings.Replace(m.className, "/", "_", -1), name)
}

func (m *Method) JavaType() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s ", m.Modifier)
	fmt.Fprintf(buf, "%s %s(", m.returnType.JavaType(), m.Name)

	for n, r := range m.arguments {
		if n > 0 {
			fmt.Fprintf(buf, ", ")
		}
		fmt.Fprintf(buf, "%s arg%d", r.JavaType(), n+1)
	}

	fmt.Fprintf(buf, ")")
	return buf.String()
}

func (m *Method) GolangType() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func %s(\n\t", m.GoFuncName())

	fmt.Fprintf(buf, "env uintptr,\n\t")

	dot := ","
	if len(m.arguments) == 0 {
		dot = ""
	}
	if m.IsStatic() {
		fmt.Fprintf(buf, "clazz uintptr%s /* java.lang.Class */", dot)
	} else {
		fmt.Fprintf(buf, "thiz uintptr%s /* %s */", dot,
			strings.Replace(m.className, "/", ".", -1))
	}

	for n, r := range m.arguments {
		dot := ","
		if n == len(m.arguments)-1 {
			dot = ""
		}
		if r.IsObject() || r.IsArray() {
			dot += " /* " + r.JavaType() + " */"
		}
		fmt.Fprintf(buf, "\n\targ%d %s%s", n+1, r.GolangType(), dot)
	}

	fmt.Fprintf(buf, ") ")
	if !m.returnType.IsVoid() {
		dot := ""
		if m.returnType.IsObject() || m.returnType.IsArray() {
			dot = "/* " + m.returnType.JavaType() + " */ "
		}
		fmt.Fprintf(buf, "%s %s", m.returnType.GolangType(), dot)
	}
	fmt.Fprintf(buf, "{")
	if !m.returnType.IsVoid() {
		fmt.Fprintf(buf, "\n\treturn %s ", m.returnType.Zero())
	}
	fmt.Fprintf(buf, "\n}")
	return buf.String()
}

func (m *Method) CDeclaration() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%s %s(GoUintptr, GoUintptr", m.returnType.CType(), m.GoFuncName())
	for _, r := range m.arguments {
		fmt.Fprintf(buf, ", %s", r.CType())
	}
	fmt.Fprintf(buf, ")")

	return buf.String()
}

func (m *Method) String() string {
	return fmt.Sprintf("//\n// Class : %s\n// Method: %s\n//\n//export %s\n%s",
		m.JavaClassName(), m.JavaType(), m.GoFuncName(), m.GolangType())
}
