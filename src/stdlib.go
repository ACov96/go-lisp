package main

import "fmt"
import "os"

var StandardLibrary, Special map[string]interface{}
func Init() {
	StandardLibrary = map[string]interface{}{}
	Special = map[string]interface{}{
		"print": func(args []Element, ctx *Context) interface{} {
			for _, el := range args {
				tmp := el.val
				if el.kind == "identifier" {
					tmp = ctx.Get(el.val.(string))
				}
				if val, ok := tmp.(string); ok {
					fmt.Printf("%s", val)
				} else if val, ok := tmp.(float64); ok {
					fmt.Printf("%f", val)
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
		// "let": func(args []interface{}, ctx *Context) interface{} {
		// 	var tuples *list.List
		// 	letContext := Context{make(map[string]interface{}), ctx}
		// 	letArgs := args[0].(ListElem)
		// 	tuples = letArgs.subTree
		// 	for el := tuples.Front(); el != nil; el = el.Next() {
		// 		var tuple *list.List
		// 		var temp ListElem
		// 		temp, ok := el.Value.(ListElem)
		// 		if !ok {
		// 			fmt.Println("Missing parens around tuple")
		// 			os.Exit(1)
		// 		}
		// 		tuple = temp.subTree
		// 		id, isID := tuple.Front().Value.(IdentifierElem)
		// 		if !isID {
		// 			fmt.Println("Left side of tuple must be identifier")
		// 			os.Exit(1)
		// 		}
		// 		val := tuple.Front().Next().Value
		// 		if alias, ok := val.(IdentifierElem); ok {
		// 			letContext.scope[id.val] = ctx.Get(alias.val)
		// 		} else {
		// 			letContext.scope[id.val] = interpret(val, &letContext)
		// 		}
		// 	}
		// 	for _, el := range args[1:] {
		// 		interpret(el, &letContext)
		// 	}
		// 	return nil
		// },
		// "lambda": func(args []interface{}, ctx *Context) interface{} {
			
		// 	return func(callArgs []interface{}, lambdaContext *Context) interface{} {
				
		// 	}
		// },
	}

}
