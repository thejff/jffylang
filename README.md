# Jffy Lang

_JFF Programming Language_

_Pronounced however you'd like; some suggestions: "Jeffy lang", "Jiffy lang"_

**Latest Version:** 0.1.0-prealpha

**DISCLAIMER: This is not a production level lanugage, use at your own risk!**

This language is written expressly following Robert Nystrom's
Book: [Crafting Interpreters](https://craftinginterpreters.com/)
I may make changes, if there is something I feel like implementing beyond what is
in the book. Those may or may not be well documented.

The book recommends following along in Java to write the interpreter, or refer 
to other people's interpretations in the language of your choice. I'm doing 
neither as I want to challenge myself. 

The interpreter and initial compiler will be written in Go, later the intention
is to (maybe) rewrite both in Jffy. 

This language is being made _Just For Fun_, it is released under the MIT Licence.

___

## Usage

There is no build provided currently. You will need to build from scratch.

Equally, there aren't any sugary things like LSPs, formatters, or the like, you
won't have any fancy colouring when writing code in Jffy just yet.

### Building
1. Clone this repo
2. Make sure you have Go installed
3. Run `make`

You now have a build of Jffy in the `./bin` folder.

### Running

#### REPL
Simply calling `./bin/jffy` will start the REPL. 

**NOTE:** The REPL needs a bit of work doing to make it nicer to work with.

#### Files
Write your code in a file, use the extensions .jff or .jffy, whichever takes your fancy.

Call `./bin/jffy yourfile.jff`

**NOTE:** Currently only a single file is supported.

## Syntax

### Variables
```go
var a;
var b = "some string";

print a; // error

a = 1;
print a; // 1; ok

```

### Printing
```go
var a = "hello";
print a; // hello
```

### Concatenation
```go
var a = "hello ";
print a .. "world"; // hello world
```

### Maths
```go
print 2 + 3; // 5
print 3 - 2; // 1
print 2 * 3; // 6
print 6 / 3; // 2

print 2 < 3; // true
print 2 > 3; // false
print 2 <= 2; // true
print 3 >= 5; // false
```

### Control Flow
```go
var a = 5;
if a == 5 {
    print "a is what it is";
} else {
    print "I should never see the light of day!";
}

var i = 0;
while i < 10 {
    print i; // 0..9
    i = i + 1;
}

for var i = 0; i < 10; i = i + 1 {
    print i; // 0..9
}
```

### Functions
```go
fun add(a, b) {
    return a + b;
}

print add(1, 2); // 3
```

#### Anonymous Functions
```go
fun add(a, b, fn) {
    fn(a + b);
}

add(2, 5, fun (total) {
    print total;
}); // 7

var a = fun (total) {
    print total * 10;
}

add(2, 3, a); // 50
```

## Documentation

[Docs](./docs/TOC.md)<br/>
[Design](./docs/DESIGN.md)<br/>
[Todo](./docs/TODO.md)<br/>
[Changelog](./docs/CHANGELOG.md)
