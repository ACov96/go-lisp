package main

import "fmt"
import "os"
import "container/list"

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
	var l []interface{} 
	l = ast
	for _, el {
		interpret(el, nil)
	}
}

func interpret(input interface{}, ctx *Context) interface{} {
	if ctx == nil {
		newContext := Context{StandardLibrary, nil}
		return interpret(input, &newContext)
	} else if el, ok := input.([]interface{}); ok {
		return interpretList(el.subTree, ctx)
	} else if el, ok := Identifier(input); ok {
		val := ctx.Get(el)
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

func interpretList(ast AST, ctx *Context) interface{} {
	var l []interface{}
	var isIdentifier, isSpecial bool = false, false
	var special interface{}
	var args []interface{}
	l = ast
	first, isIdentifier := Identifier(l[0])
	if isIdentifier {
		special, isSpecial = Special[first]
	}
	if isSpecial && len(l) > 0 {
		args = l[1:]
		return special.(func([]interface{}, *Context) interface{})(args, ctx)
	} else {
		var evaluatedList []interface{}
		for _, el := range l {
			evaluatedList = append(evaluatedList, interpret(el, ctx)))
		}
		if f, ok := evaluatedList[0].(func([]interface{}, *Context) interface{}); ok {
			args = evaluatedList[1:]
			return f(args, ctx)
		} else {
			return AST(evaluatedList)
		}	
	}
	
}
