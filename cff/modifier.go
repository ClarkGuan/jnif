package cff

import "strings"

type Modifier int

func (m Modifier) String() string {
	var prefix []string

	if m.IsPublic() {
		prefix = append(prefix, "public")
	} else if m.IsProtected() {
		prefix = append(prefix, "protected")
	} else if m.IsPrivate() {
		prefix = append(prefix, "private")
	} else {
		prefix = append(prefix, "/*package*/")
	}

	if m.IsStatic() {
		prefix = append(prefix, "static")
	}

	if m.IsFinal() {
		prefix = append(prefix, "final")
	} else if m.IsAbstract() {
		prefix = append(prefix, "abstract")
	}

	if m.IsNative() {
		prefix = append(prefix, "native")
	}

	if m.IsStrictfp() {
		prefix = append(prefix, "strictfp")
	}

	if m.IsSynchronized() {
		prefix = append(prefix, "synchronized")
	}

	if m.IsVolatile() {
		prefix = append(prefix, "volatile")
	}

	if m.IsTransient() {
		prefix = append(prefix, "transient")
	}

	return strings.Join(prefix, " ")
}

func (m Modifier) IsPublic() bool {
	return m&0x0001 == 0x0001
}

func (m Modifier) IsPrivate() bool {
	return m&0x0002 == 0x0002
}

func (m Modifier) IsProtected() bool {
	return m&0x0004 == 0x0004
}

func (m Modifier) IsStatic() bool {
	return m&0x0008 == 0x0008
}

func (m Modifier) IsFinal() bool {
	return m&0x0010 == 0x0010
}

func (m Modifier) IsVolatile() bool {
	return m&0x0040 == 0x0040
}

func (m Modifier) IsTransient() bool {
	return m&0x0080 == 0x0080
}

func (m Modifier) IsSynchronized() bool {
	return m&0x0020 == 0x0020
}

func (m Modifier) IsBridge() bool {
	return m&0x0040 == 0x0040
}

func (m Modifier) IsVarargs() bool {
	return m&0x0080 == 0x0080
}

func (m Modifier) IsNative() bool {
	return m&0x0100 == 0x0100
}

func (m Modifier) IsStrictfp() bool {
	return m&0x0800 == 0x0800
}

func (m Modifier) IsInterface() bool {
	return m&0x0200 == 0x0200
}

func (m Modifier) IsAbstract() bool {
	return m&0x0400 == 0x0400
}

func (m Modifier) IsSynthetic() bool {
	return m&0x1000 == 0x1000
}

func (m Modifier) IsAnnotation() bool {
	return m&0x2000 == 0x2000
}

func (m Modifier) IsEnum() bool {
	return m&0x4000 == 0x4000
}

func (m Modifier) IsModule() bool {
	return m&0x8000 == 0x8000
}
