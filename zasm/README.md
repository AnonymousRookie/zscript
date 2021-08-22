# zasm

zasm负责对zcomplier生成的汇编文件(.zasm)进行汇编，并生成zvm可执行文件(.zse)。

### 示例

- demo.zasm

```
; Functions
Func sum
{
	Param b
	Param a

	; return a + b;
	Push a
	Push b
	Pop _T1
	Pop _T0
	Add _T0, _T1
	Push _T0
	Pop _RetVal
	Ret 
}
Func main
{
	Var str
	Var a
	Var b
	Var c
	Var s
	Var ret

	; str = "Hello zscript!";
	Push "Hello zscript!"
	Pop _T0
	Mov str, _T0

	; a = 1.1;
	Push 1.100
	Pop _T0
	Mov a, _T0

	; b = 2;
	Push 2
	Pop _T0
	Mov b, _T0

	; c = 3;
	Push 3
	Pop _T0
	Mov c, _T0

	; s = sum(a, b);
	Push a
	Push b
	Call sum
	Push _RetVal
	Pop _T0
	Mov s, _T0

	; ret = 1 + 9 * 4 / (8 - 5) * 2 + sum(a, b) - c;
	Push 1
	Push 9
	Push 4
	Pop _T1
	Pop _T0
	Mul _T0, _T1
	Push _T0
	Push 8
	Push 5
	Pop _T1
	Pop _T0
	Sub _T0, _T1
	Push _T0
	Pop _T1
	Pop _T0
	Div _T0, _T1
	Push _T0
	Push 2
	Pop _T1
	Pop _T0
	Mul _T0, _T1
	Push _T0
	Pop _T1
	Pop _T0
	Add _T0, _T1
	Push _T0
	Push a
	Push b
	Call sum
	Push _RetVal
	Pop _T1
	Pop _T0
	Add _T0, _T1
	Push _T0
	Push c
	Pop _T1
	Pop _T0
	Sub _T0, _T1
	Push _T0
	Pop _T0
	Mov ret, _T0

	; return;
	Ret 
}
```

- 根据demo.zasm进行汇编生成demo.zse

```
./zasm demo.zasm
```

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
