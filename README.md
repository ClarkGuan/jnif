# jnif

JNI 的 Go 语言模板代码生成工具，配合 https://github.com/ClarkGuan/jni 使用，方便 Java 与 Go 之间通讯。

注：本工具用于替换 https://github.com/ClarkGuan/gojni。
gojni 主要借助 Java 语法解析器生成代码，如果解析文件缺少上下文，则很多信息拿不到。jnif 工具解析 class 文件，规避了这类问题。

#### 安装

```bash
go get github.com/ClarkGuan/jnif
```

#### 使用

工具选项说明：

* p: 生成 Go 源文件的 package 名称。默认值 "main"
* o: 生成文件所在目录的路径。默认为 $PWD 的值
* 指定的 jar 或 class 文件或包含他们的目录路径

举例：

```bash
jnif -p hello -o ../helloworld build/java/classes/HelloWorld.class
```

在目录 `../helloworld` 生成文件 `libs.c` 和 `libs.go`。
