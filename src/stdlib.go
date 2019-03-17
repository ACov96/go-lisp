package main

import "fmt"
import "os"
import "container/list"

var StandardLibrary map[string]interface{} = map[string]interface{}{
	"true": true,
	"false": false,
}

var Special map[string]interface{} = map[string]interface{}{
	"print": func(args []interface{}, ctx *Context) interface{} {
		for _, intf := range args {
			if val, ok := intf.(StringElem); ok {
				fmt.Printf("%s", val.val)
			} else if val, ok := intf.(NumberElem); ok {
				fmt.Printf("%f", val.val)
			} else if val, ok := intf.(IdentifierElem); ok {
				if (val.val == "true" || val.val == "false") {
					fmt.Printf("%t", StandardLibrary[val.val].(bool))
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
		var letArgs *list.Element
		letArgs = args[0].(*list.Element)
		argsList := letArgs.Value.(ListElem)
		PrintAST(argsList.subTree)
		
		return nil
	},
}
