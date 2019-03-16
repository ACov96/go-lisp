package main

import "container/list"
import "fmt"
import "strconv"
import "os"

type AST *list.List

type BaseElem struct {
	elemType int
}

type IdentifierElem struct {
	BaseElem
	val string
}

type StringElem struct {
	BaseElem
	val string
}

type NumberElem struct {
	BaseElem
	val float64
}

type ListElem struct {
	BaseElem
	subTree AST
}

func Parser(tokens []Token) AST {
	tokenList := list.New()
	for _, token := range(tokens) {
		tokenList.PushBack(token)
	}
	l := parenthesize(tokenList, list.New())
	return AST(l)
}

func parenthesize(tokens *list.List, tree *list.List) *list.List {
	if tokens.Len() == 0 {
		return tree
	}
	el := tokens.Remove(tokens.Front())
	currToken, ok := el.(Token)
	if !ok {
		fmt.Println("Unable to convert list elem to token")
		os.Exit(1)
	}
	if currToken.val == "(" {
		// Starting a new list
		subList := ListElem{BaseElem{LIST}, parenthesize(tokens, list.New())}
		tree.PushBack(subList)
		return parenthesize(tokens, tree)
	} else if currToken.val == ")" {
		// Ending the current list
		return tree
	} else {
		// We have some value in a list
		switch currToken.token {
		case STRING:
			
			tree.PushBack(StringElem{BaseElem{STRING}, currToken.val[1:len(currToken.val)-1]})
		case NUMBER:
			f, err := strconv.ParseFloat(currToken.val, 64)
			if err != nil {
				fmt.Println("Unable to convert number")
				os.Exit(1)
			}
			tree.PushBack(NumberElem{BaseElem{NUMBER}, f})
		default:
			tree.PushBack(IdentifierElem{BaseElem{ID}, currToken.val})
		}
		return parenthesize(tokens, tree)
	}
}

func PrintAST(ast AST) {
	walkAndPrintTree(0, ast)
}

func walkAndPrintTree(indent int, ast AST) {
	var l *list.List
	l = ast
	spacing := ""
	for i := 0; i < indent; i++ {
		spacing += "  "
	}
	for el := l.Front(); el != nil; el = el.Next() {
		// First check if list
		if listElem, ok := el.Value.(ListElem); ok {
			fmt.Printf("(\n")
			walkAndPrintTree(indent+1, listElem.subTree)
			fmt.Printf(")\n")
		} else if numElem, ok := el.Value.(NumberElem); ok {
			fmt.Printf("%s[Number: %f]\n", spacing, numElem.val)
		} else if stringElem, ok := el.Value.(StringElem); ok {
			fmt.Printf("%s[String: \"%s\"]\n", spacing, stringElem.val)
		} else {
			idElem := el.Value.(IdentifierElem)
			fmt.Printf("%s[Identifier: %s]\n", spacing, idElem.val)
		}
	}
}
