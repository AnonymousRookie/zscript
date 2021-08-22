# zvm

zvm负责加载并执行可执行文件(.zse)。


### .zse可执行文件格式

```
// 文件头
Header
// 指令流
Instruction Stream
// 字符串表
String Table
// 函数表
Function Table
```

### 执行.zse可执行文件

```
./zvm demo.zse
```
