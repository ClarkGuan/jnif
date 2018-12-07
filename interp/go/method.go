package gob

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ClarkGuan/jnif/interp"
)

func cType(jt interp.JavaType) string {
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

func goType(jt interp.JavaType) string {
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

func goZero(jt interp.JavaType) string {
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

func goFuncName(m *interp.Method) string {
	return fmt.Sprintf("%s_%s",
		strings.Replace(m.ClassName, "/", "_", -1), m.OverloadName())
}

func golangType(m *interp.Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func %s(\n\t", goFuncName(m))

	fmt.Fprintf(buf, "ev uintptr,\n\t")

	dot := ","
	if len(m.Arguments) == 0 {
		dot = ""
	}
	if m.IsStatic() {
		fmt.Fprintf(buf, "clazz uintptr%s /* java.lang.Class */", dot)
	} else {
		fmt.Fprintf(buf, "thiz uintptr%s /* %s */", dot,
			strings.Replace(m.ClassName, "/", ".", -1))
	}

	for n, r := range m.Arguments {
		dot := ","
		if n == len(m.Arguments)-1 {
			dot = ""
		}
		if r.IsObject() || r.IsArray() {
			dot += " /* " + r.JavaType() + " */"
		}
		fmt.Fprintf(buf, "\n\targ%d %s%s", n+1, goType(r), dot)
	}

	fmt.Fprintf(buf, ") ")
	if !m.ReturnType.IsVoid() {
		dot := ""
		if m.ReturnType.IsObject() || m.ReturnType.IsArray() {
			dot = "/* " + m.ReturnType.JavaType() + " */ "
		}
		fmt.Fprintf(buf, "%s %s", goType(m.ReturnType), dot)
	}
	fmt.Fprintf(buf, "{")
	if !m.ReturnType.IsVoid() {
		fmt.Fprintf(buf, "\n\treturn %s ", goZero(m.ReturnType))
	}
	fmt.Fprintf(buf, "\n}")
	return buf.String()
}

func declaration(m *interp.Method) string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%s %s(GoUintptr, GoUintptr", cType(m.ReturnType), goFuncName(m))
	for _, r := range m.Arguments {
		fmt.Fprintf(buf, ", %s", cType(r))
	}
	fmt.Fprintf(buf, ")")

	return buf.String()
}

func goFunc(m *interp.Method) string {
	return fmt.Sprintf("//\n// Class : %s\n// Method: %s\n//\n//export %s\n%s",
		m.JavaClassName(), m.JavaSignature(), goFuncName(m), golangType(m))
}
