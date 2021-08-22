# zcomplier

zcomplier负责对zscript文件(.zs)进行词法分析、语法分析，并生成zvm汇编代码文件(.zasm)。

### 特性

- 变量定义与赋值
- 函数与函数调用
- 表达式计算

### 示例

- demo.zs 

```
func sum(a, b)
{
    return a + b;
}

func main()
{
    var str;
    str = "Hello zscript!";

    var a;
    var b;
    var c;

    a = 1.1;
    b = 2;
    c = 3;
    
    var s;
    s = sum(a, b);

    var ret;
    ret = 1 + 9 * 4 / (8 - 5) * 2 + sum(a, b) - c;

    return;
}
```

- 编译demo.zs生成demo.zasm
```
./zcomplier demo.zs
```

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