package jpath_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/text/jpath"
)

func Example() {
	const json = `{
	"arranged": "alphabetically",
	"companies": [
		{"name": "apple"},
		{"name": "facebook"},
		{"name": "github"},
		{"name": "google"}
]}`
	r := func() io.Reader { return strings.NewReader(json) }
	ch, _ := jpath.Fetch(r(), ".arranged")
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Fetch(r(), ".companies[:].name")
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Fetch(r(), ".companies[1:3].name")
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Fetch(r(), ".companies[^2].name")
	for v := range ch {
		fmt.Println(v.Value)
	}

	// Output:
	// alphabetically
	// --
	// apple
	// facebook
	// github
	// google
	// --
	// facebook
	// github
	// --
	// github
}
