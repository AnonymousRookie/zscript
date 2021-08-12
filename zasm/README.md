# zasm

zasm负责对zcomplier生成的汇编文件(.zasm)进行汇编，并生成zvm可执行文件(.zse)。


### .zse可执行文件格式

```
Header

Instruction Stream

String Table

Function Table
```