package ir

const headerTpl = `// 此文件为动态生成的，请不要修改！

#include <stdlib.h>
#include <stdio.h>
#include <stddef.h>
#include <stdint.h>
#include <jni.h>

#include "_cgo_export.h"

extern void jniOnLoad(GoUintptr vm);
extern void jniOnUnload(GoUintptr vm);

{{range $key, $val := .classes}}// Class: {{$key}}
{{range $val}}extern {{.CDeclaration}};
{{end}}
{{end}}jint JNI_OnLoad(JavaVM *vm, void *reserved) {
    JNIEnv *env = NULL;
    if ((*vm)->GetEnv(vm, (void **) &env, JNI_VERSION_1_6) != JNI_OK) {
        fprintf(stderr, "[%s:%d] GetEnv() return error\n", __FILE__, __LINE__);
        abort();
    }

    jclass clazz;
    JNINativeMethod methods[{{.maxCount}}];
    jint size;
    char *name;

    {{range $key, $val := .classes}}name = "{{$key}}";
    clazz = (*env)->FindClass(env, name);
    if (clazz == NULL) {
        fprintf(stderr, "[%s:%d] FindClass() \"%s\" return error\n", __FILE__, __LINE__, name);
        abort();
    }
    size = 0;

    {{range $val}}methods[size].fnPtr = {{.GoFuncName}};
    methods[size].name = "{{.Name}}";
    methods[size].signature = "{{.Desc}}";
    size++;

    {{end}}if ((*env)->RegisterNatives(env, clazz, methods, size) != 0) {
        fprintf(stderr, "[%s:%d] %s RegisterNatives() return error\n", __FILE__, __LINE__, name);
        abort();
    }

    {{end}}jniOnLoad((GoUintptr) vm);
    return JNI_VERSION_1_6;
}

void JNI_OnUnload(JavaVM *vm, void *reserved) {
    jniOnUnload((GoUintptr) vm);
}
`

const goTpl = `package {{.packageName}}

//
// #include <stdlib.h>
// #include <stddef.h>
// #include <stdint.h>
import "C"

//export jniOnLoad
func jniOnLoad(vm uintptr) {
    // TODO
}

//export jniOnUnload
func jniOnUnload(vm uintptr) {
    // TODO
}

{{range $key, $val := .classes}}// Class: {{$key}}
{{range $val}}{{.}}

{{end}}{{end}}
`
