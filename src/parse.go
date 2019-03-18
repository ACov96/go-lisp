package main

import "container/list"
import "fmt"
import "strconv"
import "os"

type AST []interface{}
type Identifier string

func Parser(tokens []Token) AST {
	tokenList := list.New()
	for _, token := range(tokens) {
		tokenList.PushBack(token)
	}
	l := parenthesize(tokenList, list.New())
	return AST(makeListArray(l))
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
		subList := parenthesize(tokens, list.New())
		tree.PushBack(subList)
		return parenthesize(tokens, tree)
	} else if currToken.val == ")" {
		// Ending the current list
		return tree
	} else {
		// We have some value in a list
		switch currToken.token {
		case STRING:
			tree.PushBack(currToken.val[1:len(currToken.val)-1])
		case NUMBER:
			f, err := strconv.ParseFloat(currToken.val, 64)
			if err != nil {
				fmt.Println("Unable to convert number")
				os.Exit(1)
			}
			tree.PushBack(f)
		case BOOL:
			b, err := strconv.ParseBool(currToken.val)
			if err != nil {
				fmt.Println("Unable to convert boolean")
				os.Exit(1)
			}
			tree.PushBack(b)
		default:
			tree.PushBack(Identifier(currToken.val))
		}
		return parenthesize(tokens, tree)
	}
}

func makeListArray(l *list.List) []interface{} {
	var arr []interface{}
	for el := l.Front(); el != nil; el = el.Next() {
		if subList, ok := el.Value.(*list.List); ok {
			arr = append(arr, makeListArray(subList))
		} else {
			arr = append(arr, el.Value)
		}
	}
	return arr
}

func PrintAST(ast AST) {
	walkAndPrintTree(0, ast)
}

func walkAndPrintTree(indent int, ast AST) {
	var l []interface{}
	l = ast
	spacing := ""
	for i := 0; i < indent; i++ {
		spacing += "  "
	}
	for _, el := range l {
		if listElem, ok := el.([]interface{}); ok {
			fmt.Printf("%s(\n", spacing)
			walkAndPrintTree(indent+1, listElem)
			fmt.Printf("%s)\n", spacing)
		} else if numElem, ok := el.(float64); ok {
			fmt.Printf("%s[Number: %f]\n", spacing, numElem)
		} else if idElem, ok := el.(Identifier); ok {
			fmt.Printf("%s[Identifier: %s]\n", spacing, idElem)
		} else if boolElem, ok := el.(bool); ok {
			fmt.Printf("%s[Boolean: %t]\n", spacing, boolElem)
		} else {
			stringElem := el.(string)
			fmt.Printf("%s[String: \"%s\"]\n", spacing, stringElem)
		}
	}
}
