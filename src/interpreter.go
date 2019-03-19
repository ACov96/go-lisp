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
	} else if el, ok := input.(Element); ok && el.kind == "identifier"{
		val := ctx.Get(el.val.(string))
		if val != nil {
			return val
		} else {
			fmt.Printf("Cannot find identifier: %s", el)
			os.Exit(1)
		}
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
