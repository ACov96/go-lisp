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

type Function func([]interface{}) interface{}

func Interpret(ast AST) {
	var l *list.List
	l = ast
	for el := l.Front(); el != nil; el = el.Next() {
		interpret(el.Value, nil)
	}
}

func interpret(input interface{}, ctx *Context) interface{} {
	if ctx == nil {
		newContext := Context{StandardLibrary, nil}
		return interpret(input, &newContext)
	} else if el, ok := input.(ListElem); ok {
		return interpretList(el.subTree, ctx)
	} else if el, ok := input.(IdentifierElem); ok {
		val := ctx.Get(el.val)
		if val != nil {
			return val
		} else {
			fmt.Printf("Cannot find identifier: %s", el.val)
			os.Exit(1)
		}
	} else if el, ok := input.(NumberElem); ok {
		return el.val
	} else {
		el, ok := input.(StringElem)
		if !ok {
			fmt.Println("Unable to resolve value")
			os.Exit(1)
		}
		return el.val
	}
	return nil
}

func interpretList(ast AST, ctx *Context) interface{} {
	var l *list.List
	l = ast
	evaluatedList := list.New()
	for el := l.Front(); el != nil; el = el.Next() {
		evaluatedList.PushBack(interpret(el.Value, ctx))
	}
	if f, ok := evaluatedList.Front().Value.(Function); ok {
		var args []interface{}
		for el := evaluatedList.Front().Next(); el != nil; el = el.Next() {
			args = append(args, el)
		}
		return f(args)
	} else {
		return ListElem{BaseElem{LIST}, AST(evaluatedList)}
	}
}
