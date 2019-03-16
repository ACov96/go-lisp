package main

import "fmt"
import "os"
import "container/list"

var StandardLibrary map[string]interface{} = map[string]interface{}{
	"print": func(args []interface{}) interface{} {
		for idx, intf := range args {
			var el *list.Element
			el = intf.(*list.Element)
			if val, ok := el.Value.(string); ok {
				fmt.Printf(val)
			} else if val, ok := el.Value.(float64); ok {
				fmt.Printf("%d", val)
			} else {
				fmt.Printf("Unable to resolve interface: %v\n", el.Value.(StringElem))
				os.Exit(1)
			}
			if idx != len(args) {
				fmt.Printf(" ")
			}
			fmt.Printf("\n")
		} 
		return nil
	},
}
