package main

import "container/list"
import "fmt"
import "strconv"
import "os"

type Element struct {
	kind string
	val interface{}
}
type AST []Element


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
		tree.PushBack(Element{"list", subList})
		return parenthesize(tokens, tree)
	} else if currToken.val == ")" {
		// Ending the current list
		return tree
	} else {
		// We have some value in a list
		switch currToken.token {
		case STRING:
			tree.PushBack(Element{"literal", currToken.val[1:len(currToken.val)-1]})
		case NUMBER:
			f, err := strconv.ParseFloat(currToken.val, 64)
			if err != nil {
				fmt.Println("Unable to convert number")
				os.Exit(1)
			}
			tree.PushBack(Element{"literal", f})
		case BOOL:
			b, err := strconv.ParseBool(currToken.val)
			if err != nil {
				fmt.Println("Unable to convert boolean")
				os.Exit(1)
			}
			tree.PushBack(Element{"literal", b})
		default:
			tree.PushBack(Element{"identifier", currToken.val})
		}
		return parenthesize(tokens, tree)
	}
}

func makeListArray(l *list.List) []Element {
	var arr []Element
	for el := l.Front(); el != nil; el = el.Next() {
		if elType := el.Value.(Element); elType.kind == "list" {
			arr = append(arr, Element{"list", makeListArray(elType.val.(*list.List))})
		} else {
			arr = append(arr, el.Value.(Element))
		}
	}
	return arr
}

func PrintAST(ast AST) {
	walkAndPrintTree(0, ast)
}

func walkAndPrintTree(indent int, ast AST) {
	var l []Element
	l = ast
	spacing := ""
	for i := 0; i < indent; i++ {
		spacing += "  "
	}
	for _, el := range l {
		if el.kind == "list" {
			fmt.Printf("%s(\n", spacing)
			walkAndPrintTree(indent+1, el.val.([]Element))
			fmt.Printf("%s)\n", spacing)
		} else if el.kind == "identifier" {
			fmt.Printf("%s[Identifier: %s]\n", spacing, el.val.(string))
		} else if val, ok := el.val.(float64); ok {
			fmt.Printf("%s[Number: %f]\n", spacing, val)
		} else if val, ok := el.val.(bool); ok {
			fmt.Printf("%s[Boolean: %t]\n", spacing, val)
		} else if val, ok := el.val.(string); ok {
			fmt.Printf("%s[String: \"%s\"]\n", spacing, val)
		} else {
			fmt.Printf("Unknown\n")
		}
	}
}
