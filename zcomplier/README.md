# zcomplier

zcomplier负责对zscript文件(.zs)进行词法分析、语法分析，并生成zvm汇编代码文件(.zasm)。

### 特性

- 变量定义与赋值

```
var a;
a = 1;

var b;
b = 2.0;

var c;
c = "hello zscript";
```

- 函数与函数调用

```
func sum(a, b)
{
    return a + b;
}

func run()
{
    var a;
    a = 1;

    var b;
    b = 2.0;
    
    var s;
    s = sum(a, b);
}
```

- 表达式

```
func sum(a, b)
{
    return a + b;
}

func main()
{
    var a;
    var b;
    var c;

    a = 1.1;
    b = 2;
    c = 3;
    
    var ret;
    ret = 1 + 9 * 4 / (8 - 5) * 2 + sum(a, b) - c;
}
```
