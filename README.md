# Golox

Implementation of the Lox language in Golang (https://craftinginterpreters.com)

## Install

```bash
go install github.com/taki-mekhalfa/golox@latest
```

## Usage

### Prompt

```bash
golox
```

#### Examples

```bash
>> var a = 3;
>> var b = 5;
>> print a + b + 1;
9
>> var space = " ";          
>> print "Hello" + space + "world!";
Hello world!
>> 
```

### Intrepret a file

```bash
golox src.lox
```

#### Examples

```c
// fib returns the n-th Fibonacci element
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1) {
  print fib(i);
}

/*
Output:
0
1
1
2
3
5
8
13
21
34
/*
```

```c
// makeCounter demonstrates the use of closures
fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print i;
  }

  return count;
}

var counter = makeCounter();
counter(); // "1".
counter(); // "2".

/*
Output:
1
2
/*
```

```c
fun count(n) {
  while (n < 100) {
    if (n == 3) return n; // <-- return whenever you want
    print n;
    n = n + 1;
  }
}

count(1);

/*
Output:
1
2
/*
```

```c
// nest blocks however you want
{
    {
        var i = 0;
        while (i < 3) {
            print i;
            i = i + 1;
        }
    }
}

/*
Output:
0
1
2
/*
```

```c
// manage scopes with blocks
var a = "global a";
var b = "global b";
var c = "global c";
{
  var a = "outer a";
  var b = "outer b";
  {
    var a = "inner a";
    print a;
    print b;
    print c;
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;

/*
Output:
inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c
/*
```

```c
// `print` and `clock` are built-in functions
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

var t1 = clock();
for (var i = 0; i < 30; i = i + 1) {
    fib(i);
}
print "running time(s):";
print clock() - t1; // not very efficient haha :)
```
