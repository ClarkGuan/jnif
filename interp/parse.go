package interp

import (
	"fmt"

	"github.com/ClarkGuan/jnif/cff"
)

func Print(args []string) error {
	if len(args) == 0 || gRegisters[args[0]] == nil {
		return noRegisterErr
	}

	transformer := gRegisters[args[0]]
	if path, err := transformer.Init(args[1:]); err != nil {
		return err
	} else if methods, max, err := parse(path); err != nil {
		return err
	} else if err = transformer.Transform(methods, max); err != nil {
		return err
	}

	return nil
}

func parse(path string) (methods map[string][]*Method, max int, err error) {
	if infos, err := cff.Parse(path); err != nil {
		return nil, 0, err
	} else {
		if methods == nil {
			methods = make(map[string][]*Method)
		}

		for _, info := range infos {
			funcs := make(map[string]int)
			for _, m := range info.Methods {
				if !m.AccessFlags.IsNative() {
					continue
				}

				method := Method{}
				method.Name = m.Name
				method.ClassName = info.Name
				method.Desc = m.Desc
				method.Arguments, method.ReturnType = parseDesc(m.Desc)
				method.Modifier = m.AccessFlags

				if count, ok := funcs[m.Name]; ok {
					method.overload = count + 1
					funcs[m.Name] = count + 1
				} else {
					funcs[m.Name] = 0
				}

				methods[info.Name] = append(methods[info.Name], &method)
			}

			size := len(methods[info.Name])
			if size > max {
				max = size
			}
		}
	}

	return
}

type mode uint

func (m mode) setArgumentRange() mode {
	return m &^ modeRangeMusk
}

func (m mode) setReturnRange() mode {
	return m | modeReturn
}

func (m mode) setNormalState() mode {
	return m &^ modeStateMask
}

func (m mode) setArrayState() mode {
	return m | modeArray
}

func (m mode) setObjectState() mode {
	return m | modeObject
}

const (
	modeArgument  mode = 0x0 // 处理方法参数
	modeReturn    mode = 0x1 // 处理方法返回值
	modeRangeMusk      = modeArgument | modeReturn

	modeNormal    mode = 0x0 // java 基本类型
	modeArray     mode = 0x2 // java 数组
	modeObject    mode = 0x4 // java 对象
	modeStateMask      = modeNormal | modeArray | modeObject
)

func parseDesc(desc string) (arguments []JavaType, returnType JavaType) {
	mode := modeNormal
	last := -1
	for i, b := range desc {
		switch b {
		case '(':
			mode = mode.setArgumentRange()

		case ')':
			mode = mode.setReturnRange()

		case 'L':
			if mode&modeObject != modeObject {
				mode = mode.setObjectState()
				if last == -1 {
					last = i
				}
			}

		case '[':
			if mode&modeArray != modeArray {
				mode = mode.setArrayState()
				if last == -1 {
					last = i
				}
			}

		case 'V', 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', ';':
			if mode&modeObject == modeObject && b != ';' {
				continue
			}

			start := i
			if mode&^modeRangeMusk != modeNormal {
				if last == -1 && b == ';' {
					panic(fmt.Errorf("error state = %s", desc))
				}
				start = last
			}

			if mode&modeReturn == modeReturn {
				returnType = JavaType(desc[start : i+1])
				return
			} else {
				arguments = append(arguments, JavaType(desc[start:i+1]))
			}

			mode &^= modeStateMask
			last = -1
		}
	}

	panic(fmt.Errorf("can't reach here = %s", desc))
}
