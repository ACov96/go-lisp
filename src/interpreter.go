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
	} else if el, ok := input.(BoolElem); ok {
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
	var isIdentifier, isSpecial bool = false, false
	var special interface{}
	var args []interface{}
	l = ast
	first, isIdentifier := l.Front().Value.(IdentifierElem)
	if isIdentifier {
		special, isSpecial = Special[first.val]
	}
	if isSpecial && l.Len() > 0 {
		for el := l.Front().Next(); el != nil; el = el.Next() {
			args = append(args, el.Value)
		}
		return special.(func([]interface{}, *Context) interface{})(args, ctx)
	} else {
		evaluatedList := list.New()
		for el := l.Front(); el != nil; el = el.Next() {
			evaluatedList.PushBack(interpret(el.Value, ctx))
		}
		if f, ok := evaluatedList.Front().Value.(func([]interface{}, *Context) interface{}); ok {
			for el := evaluatedList.Front().Next(); el != nil; el = el.Next() {
				args = append(args, el.Value)
			}
			return f(args, ctx)
		} else {
			return ListElem{BaseElem{LIST}, AST(evaluatedList)}
		}	
	}
	
}
