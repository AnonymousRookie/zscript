# zscript

### zcomplier

zscript编译器：负责对zscript文件(.zs)进行词法分析、语法分析，并生成zvm汇编代码文件(.zasm)。

### zasm

zscript汇编器：负责对zcomplier生成的汇编文件(.zasm)进行汇编，并生成zvm可执行文件(.zse)。

### zvm

zscript虚拟机：负责加载并执行可执行文件(.zse)。

