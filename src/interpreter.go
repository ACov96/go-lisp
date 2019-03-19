package main

import "fmt"
import "os"

type Context struct {
	scope map[string]interface{}
	parent *Context
}

func (c Context) Get(x string) interface{} {
	if val, ok := c.scope[x]; ok {
		return val
	} else if c.parent != nil {
		return c.parent.Get(x)
	} else {
		return nil
	}
}

var StandardLibrary, Special map[string]interface{}

func Interpret(ast AST) {
	Init()
	var l []Element
	l = ast
	for _, el := range l {
		interpret(el, nil)
	}
}

func interpret(input interface{}, ctx *Context) interface{} {
	if ctx == nil {
		newContext := Context{StandardLibrary, nil}
		return interpret(input, &newContext)
	} else if el, ok := input.(Element); ok && el.kind == "list" {
		return interpretList(el.val.([]Element), ctx)
	} else if el, ok := input.(Element); ok && el.kind == "identifier" {
		val := ctx.Get(el.val.(string))
		if val != nil {
			return val
		} else {
			fmt.Printf("Cannot find identifier: %s", el)
			os.Exit(1)
		}
	} else if el, ok := input.(Element); ok && el.kind == "literal" {
		return el.val
	} else if el, ok := input.(float64); ok {
		return el
	} else if el, ok := input.(bool); ok {
		return el
	} else {
		el, ok := input.(string)
		if !ok {
			fmt.Println("Unable to resolve value")
			os.Exit(1)
		}
		return el
	}
	return nil
}

func interpretList(l []Element, ctx *Context) interface{} {
	var args []interface{}
	var isSpecial bool = false
	var special interface{}
	if l[0].kind == "identifier" {
		special, isSpecial = Special[l[0].val.(string)]
	}
	if isSpecial && len(l) > 0 {
		return special.(func([]Element, *Context) interface{})(l[1:], ctx)
	} else {
		var evaluatedList []interface{}
		for _, el := range l {
			evaluatedList = append(evaluatedList, interpret(el, ctx))
		}
		if f, ok := evaluatedList[0].(func([]interface{}, *Context) interface{}); ok {
			args = evaluatedList[1:]
			return f(args, ctx)
		} else {
			return evaluatedList
		}	
	}
	
}

func printList(list []Element, ctx *Context) {
	fmt.Printf("(")
	for idx, el := range list {
		var tmp interface{}
		if el.kind == "list" {
			printList(el.val.([]Element), ctx)
		} else if el.kind == "identifier" {
			tmp = ctx.Get(el.val.(string))
		} else {
			tmp = el.val
		}

		if val, ok := tmp.(string); ok {
			fmt.Printf("%s", val)
		} else if val, ok := tmp.(float64); ok {
			fmt.Printf("%g", val)
		} else if val, ok := tmp.(bool); ok {
			fmt.Printf("%t", val)
		}
		if idx != len(list)-1 {
			fmt.Printf(" ")
		}
	}
	fmt.Printf(")")
}

func printInterfaceList(list []interface{}) {
	fmt.Printf("(")
	for idx, el := range list {
		if val, ok := el.(string); ok {
			fmt.Printf("%s", val)
		} else if val, ok := el.(float64); ok {
			fmt.Printf("%g", val)
		} else {
			fmt.Printf("%t", el.(bool))
		}
		if idx != len(list)-1 {
			fmt.Printf(" ")
		}
	}
	fmt.Printf(")\n")
}

func Init() {
	StandardLibrary = map[string]interface{}{}
	Special = map[string]interface{}{
		"print": func(args []Element, ctx *Context) interface{} {
			for _, el := range args {
				tmp := el.val
				if el.kind == "identifier" {
					tmp = ctx.Get(el.val.(string))
					// This check is needed in case a list has already been evaluated
					if list, ok := tmp.([]interface{}); ok {
						printInterfaceList(list)
						continue
					}
				} else if el.kind == "list" {
					printList(el.val.([]Element), ctx)
					continue
				}
				if val, ok := tmp.(string); ok {
					fmt.Printf("%s", val)
				} else if val, ok := tmp.(float64); ok {
					fmt.Printf("%g", val)
				} else if val, ok := tmp.(bool); ok {
					fmt.Printf("%t", val);
				} else {
					fmt.Println("Unknown value")
					os.Exit(1)
				}
			} 
			fmt.Printf("\n")
			return nil
		},
		"let": func(args []Element, ctx *Context) interface{} {
			tuples := args[0]
			if tuples.kind != "list" {
				fmt.Println("First argument to let is not a list of tuples")
				os.Exit(1)
			}
			newContext := Context{make(map[string]interface{}), ctx}
			for _, tuple := range tuples.val.([]Element) {
				left := tuple.val.([]Element)[0]
				right := tuple.val.([]Element)[1]
				if left.kind != "identifier" {
					fmt.Println("Left side of tuple is not identifier")
					os.Exit(1)
				}
				id := left.val.(string)
				newContext.scope[id] = interpret(right, ctx)
			}
			return interpret(args[1], &newContext)
		},
	}
}
