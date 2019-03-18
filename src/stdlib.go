package main

import "fmt"
import "os"
import "container/list"

var StandardLibrary, Special map[string]interface{}
func Init() {
	StandardLibrary = map[string]interface{}{}
	Special = map[string]interface{}{
		"print": func(args []interface{}, ctx *Context) interface{} {
			for _, intf := range args {
				if val, ok := intf.(StringElem); ok {
					fmt.Printf("%s", val.val)
				} else if val, ok := intf.(NumberElem); ok {
					fmt.Printf("%f", val.val)
				} else if val, ok := intf.(BoolElem); ok {
					fmt.Printf("%t", val.val)
				} else if val, ok := intf.(IdentifierElem); ok {
					temp := ctx.Get(val.val)
					if str, isString := temp.(StringElem); isString {
						fmt.Printf("%s", str.val)
					} else if number, isNumber := temp.(NumberElem); isNumber {
						fmt.Printf("%f", number.val)
					} else if boolean, isBoolean := temp.(BoolElem); isBoolean {
						fmt.Printf("%t", boolean.val)
					} else {
						fmt.Println("Unrecognized element in context")
						os.Exit(1)
					}	
				} else {
					fmt.Printf("Unable to resolve interface: %v\n", intf)
					os.Exit(1)
				}
			} 
			fmt.Printf("\n")
			return nil
		},
		"let": func(args []interface{}, ctx *Context) interface{} {
			var tuples *list.List
			letContext := Context{make(map[string]interface{}), ctx}
			letArgs := args[0].(ListElem)
			tuples = letArgs.subTree
			for el := tuples.Front(); el != nil; el = el.Next() {
				var tuple *list.List
				var temp ListElem
				temp, ok := el.Value.(ListElem)
				if !ok {
					fmt.Println("Missing parens around tuple")
					os.Exit(1)
				}
				tuple = temp.subTree
				id, isID := tuple.Front().Value.(IdentifierElem)
				if !isID {
					fmt.Println("Left side of tuple must be identifier")
					os.Exit(1)
				}
				val := tuple.Front().Next().Value
				if alias, ok := val.(IdentifierElem); ok {
					letContext.scope[id.val] = ctx.Get(alias.val)
				} else {
					letContext.scope[id.val] = interpret(val, &letContext)
				}
			}
			for _, el := range args[1:] {
				interpret(el, &letContext)
			}
			return nil
		},
		"lambda": func(args []interface{}, ctx *Context) interface{} {
			
			return func(callArgs []interface{}, lambdaContext *Context) interface{} {
				
			}
		},
	}

}
