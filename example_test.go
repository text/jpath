package jpath_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/text/jpath"
)

func ExampleEvaluate() {
	const json = `{
	"arranged": "alphabetically",
	"companies": [
		{"name": "apple"},
		{"name": "facebook"},
		{"name": "github"},
		{"name": "google"}
]}`
	r := func() io.Reader { return strings.NewReader(json) }
	ch, _ := jpath.Evaluate(".arranged", r())
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Evaluate(".companies[:].name", r())
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Evaluate(".companies[1:3].name", r())
	for v := range ch {
		fmt.Println(v.Value)
	}
	fmt.Println("--")
	ch, _ = jpath.Evaluate(".companies[^2].name", r())
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
