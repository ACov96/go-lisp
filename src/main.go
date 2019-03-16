package main

import "os"

func main() {
	args := os.Args[1:]
	tokens := LexFile(args[0])
	ast := Parser(tokens)
	PrintAST(ast)
}
