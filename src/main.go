package main

import "os"

func main() {
	args := os.Args[1:]
	LexFile(args[0])
}