func sum(a, b)
{
    return a + b;
}

func main()
{
    var str;
    str = "Hello zscript!";
    print(str);

    var a;
    var b;
    var c;

    a = 1.1;
    b = 2;
    c = 3;
    
    var s;
    s = sum(a, b);
    print(s);

    var ret;
    ret = 1 + 9 * 4 / (8 - 5) * 2 + sum(a, b) - c;
    print(ret);

    return;
}