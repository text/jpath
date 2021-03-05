package jpath_test

import (
	"fmt"
	"strings"

	"github.com/text/jpath"
)

func ExampleFetch() {
	r := strings.NewReader(`[{"company": {"name": "apple"}},
	{"company": {"name": "facebook"}},
	{"company": {"name": "github"}},
	{"company": {"name": "google"}}]`)
	ch, _ := jpath.Fetch(r, "[*].company.name")
	for a := range ch {
		fmt.Println(a.Value)
	}
	// Output:
	// apple
	// facebook
	// github
	// google
}
