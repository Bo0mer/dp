/*
	Package flagvar extends the standard flag package with some well-known types,
	that are not available there.

	Example usage:

		package main

		import (
			"flag"
			"fmt"

			"github.com/Bo0mer/flagvar"
		)

		var (
			tags       flagvar.Array
			attributes flagvar.Map
		)

		func init() {
			flag.Var(&tags, "tag", "Tag to add.")
			flag.Var(&attributes, "attribute", "Attribute to add.")
		}

		func main() {
			flag.Parse()

			fmt.Println("tags provided:", tags)
			fmt.Println("attributes provided:", attributes)

			// When started with:
			// -tag=1 -tag=2 -attribute foo:bar -attribute baz:boo
			// Outputs:
			// tags provided: [1 2]
			// attributes provided: map[foo:bar baz:boo]
		}
*/
package flagvar
