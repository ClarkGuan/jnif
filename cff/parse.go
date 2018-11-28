package cff

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FieldInfo struct {
	Name        string
	Desc        string
	AccessFlags Modifier
}

func (info *FieldInfo) String() string {
	return fmt.Sprintf("%s %s %s", info.AccessFlags, info.Desc, info.Name)
}

type MethodInfo struct {
	Name        string
	Desc        string
	AccessFlags Modifier
}

func (info *MethodInfo) String() string {
	return fmt.Sprintf("%s %s%s", info.AccessFlags, info.Name, info.Desc)
}

type ClassInfo struct {
	Name        string
	AccessFlags Modifier
	SuperClass  string
	Interfaces  []string
	Fields      []*FieldInfo
	Methods     []*MethodInfo
}

func (info *ClassInfo) String() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%s %s ", info.AccessFlags, info.Name)

	if len(info.SuperClass) > 0 {
		fmt.Fprintf(buf, "extends %s ", info.SuperClass)
	}

	if len(info.Interfaces) > 0 {
		fmt.Fprintf(buf, "implements ")
		for n, s := range info.Interfaces {
			if n != 0 {
				fmt.Fprintf(buf, ", ")
			}
			fmt.Fprintf(buf, "%s", s)
		}
		fmt.Fprintf(buf, " ")
	}
	fmt.Fprintln(buf, "{")

	for _, f := range info.Fields {
		fmt.Fprintf(buf, "    %s;\n", f)
	}

	for _, m := range info.Methods {
		fmt.Fprintf(buf, "    %s {}\n", m)
	}

	fmt.Fprintln(buf, "}")
	return buf.String()
}

func Parse(path string) (infos []*ClassInfo, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".class") {
			if data, err := ioutil.ReadFile(path); err != nil {
				return err
			} else if classInfo, err := ParseData(data); err != nil {
				return err
			} else {
				infos = append(infos, classInfo)
			}
		} else if strings.HasSuffix(info.Name(), ".jar") {
			jarFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer jarFile.Close()

			jarInfo, err := jarFile.Stat()
			if err != nil {
				return err
			}

			zipReader, err := zip.NewReader(jarFile, jarInfo.Size())
			if err != nil {
				return err
			}

			for _, f := range zipReader.File {
				if !strings.HasSuffix(f.Name, ".class") {
					continue
				}

				//fmt.Println("parse", path, "@", f.Name)

				fc, err := f.Open()
				if err != nil {
					return err
				}
				data, err := ioutil.ReadAll(fc)
				if err != nil {
					fc.Close()
					return err
				}
				classInfo, err := ParseData(data)
				if err != nil {
					fc.Close()
					return err
				}
				infos = append(infos, classInfo)
				fc.Close()
			}

		}

		return nil
	})

	return
}

func ParseData(data []byte) (info *ClassInfo, err error) {
	cfi, err := parseClassFile(data)
	if err != nil {
		return nil, err
	}

	info = new(ClassInfo)
	info.Name = cfi.Utf8Info(cfi.constantPool[cfi.thisClass].(*constantClassInfo).nameIndex).String()
	info.AccessFlags = Modifier(cfi.accessFlags)
	if cfi.superClass != 0 {
		info.SuperClass = cfi.Utf8Info(cfi.constantPool[cfi.superClass].(*constantClassInfo).nameIndex).String()
	}
	for _, i := range cfi.interfaces {
		info.Interfaces = append(info.Interfaces, cfi.Utf8Info(cfi.constantPool[i].(*constantClassInfo).nameIndex).String())
	}
	for _, f := range cfi.fields {
		mi := FieldInfo{}
		mi.AccessFlags = Modifier(f.accessFlags)
		mi.Name = cfi.Utf8Info(f.nameIndex).String()
		mi.Desc = cfi.Utf8Info(f.descriptorIndex).String()
		info.Fields = append(info.Fields, &mi)
	}
	for _, m := range cfi.methods {
		mi := MethodInfo{}
		mi.AccessFlags = Modifier(m.accessFlags)
		mi.Name = cfi.Utf8Info(m.nameIndex).String()
		mi.Desc = cfi.Utf8Info(m.descriptorIndex).String()
		info.Methods = append(info.Methods, &mi)
	}

	return
}
