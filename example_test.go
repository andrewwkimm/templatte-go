package templatte_test

import (
	"fmt"

	"github.com/andrewwkimm/templatte-go"
)

func ExampleGreet() {
	greeting, err := templatte.Greet("Gopher")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(greeting.Message)
	// Output: Hello, Gopher!
}
