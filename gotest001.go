package main

import (
	. "github.com/pyrocto/gotest001/foo"
	"fmt"
)

func main() {
	var f Foo[int] = Bar[int]{5}
	fmt.Println("hello")
}
