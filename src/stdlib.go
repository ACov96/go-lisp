package main

import "fmt"
import "os"
import "container/list"

var StandardLibrary map[string]interface{} = map[string]interface{}{
	"print": func(args []interface{}, ctx *Context) interface{} {
		for _, intf := range args {
			var el *list.Element
			el = intf.(*list.Element)
			if val, ok := el.Value.(string); ok {
				fmt.Printf("%s", val)
			} else if val, ok := el.Value.(float64); ok {
				fmt.Printf("%f", val)
			} else if val, ok := el.Value.(bool); ok {
				fmt.Printf("%t", val)
			} else {
				fmt.Printf("Unable to resolve interface: %v\n", el.Value.(StringElem))
				os.Exit(1)
			}
		} 
		fmt.Printf("\n")
		return nil
	},
	"T": true,
	"F": false,
}
