# Golox

Implementation of a lightweight interpreter for the small dynamically typed Lox language in Golang (https://craftinginterpreters.com)

## Install

```bash
go install github.com/taki-mekhalfa/golox@latest
```

## Architecture

* A hand-crafted scanner, a stream of characters in, a stream of tokens out with error messages and error lines
* A hand-made recursive descent parser to transform the stream of tokens into an `AST` to be interpreted in a later stages. Operators precedence is built into the grammar by recursively parsing operators with high precedence before those with lower precedence. e.g.:

```c
term   → factor ( "+" factor )* ;
factor → NUMBER ( "*" NUMBER )* ;
```
* A visitor printer that pretty prints the `AST` to check that parsing is correct
* A visitor resolver, that makes a pass through the `AST` before interpretation to resolve variables binding and check for some semantic errors (returns outside a function, declared but not used variables, used but non declared variables, reference to `this` outside a method, etc.)
* A visitor tree-walk interpreter that walks through the `AST` to interpret the program. 

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
class Cake {
  init(flavor) {
    this.flavor = flavor;
  }

  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}

var cake = Cake("German chocolate");
cake.taste(); // Prints "The German chocolate cake is delicious!".
/*
Output:
The German chocolate cake is delicious!
*/
```

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
var a = "global";
{
  fun showA() {
    // close on the the global `a`
    print a;
  }

  showA();
  // redeclare `a` on the same scope 
  // on which `showA` is closing on
  var a = "block";
  showA();
  print a;

/*
Output:
global
global
block
/*
}
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

```c
// you can't redeclare the same variable in the same local scope
{
    var a = 3;
    var a = 4;
    print a;
}
/*
Output:
[line 3] Runtime Error: Already a variable with this name in this scope.
*/

// you can't return from a top level code
{
    var a = 3;
    print a;
    return a;
}
/*
Output:
[line 4] Runtime Error: Can't return from top-level code.
*/

// you should use a variable you declare
{
    var a = 3;
    var b = 3;
    var c;
}
/*
Output:
[line 3] Runtime Error: b declared but not used.
[line 4] Runtime Error: c declared but not used.
*/

// you can't return a value inside a constructor
class Klass {
  init() {
    this.age = 10;
    return "hello world";
  }
}

Klass();
/*
Output:
[line 4] Runtime Error: Can't return a value from class initializer.
*/
```
