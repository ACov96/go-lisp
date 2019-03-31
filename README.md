# Go Lisp

This is my implementation of a small, toy Lisp language written in Go. It is based off of [Mary Rose Cook's Little Lisp](https://maryrosecook.com/blog/post/little-lisp-interpreter).

## Build

You will need a current version of Go (I developed using 1.12). You will also need to install [Nex](https://github.com/blynn/nex)Then run:

```bash
make build
```

This will create an executable called `gl`, which is the interpreter.

## Usage

Pass the file name to the interpreter to run it:

```bash
gl hello.gl
```

## Documentation

There isn't any official documentation since this is just a proof of concept project. But here's what's been implemented:

**print**
```lisp
(print exp1 exp2 exp3)
```

This takes any number of arguments and prints them to stdout.

**let**
```lisp
(let ((h 'hello')
      (w 'world'))
      (print h ' ' w))
```

Create a scope and declare variables within that scope.

**if**
```lisp
(if exp
    trueExp
    falseExp)
```

Checks the boolean value of *exp* and executes *trueExp* if *exp* is true, otherwise it executes *falseExp*.

**lambda**
```lisp
(lambda (x) (print x))
```

My lambda implementation is still pretty rough, but you can define small functions with ease.
