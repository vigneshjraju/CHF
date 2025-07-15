package main

import "fmt"

func main() {

	var p *int

	i := 42
	p = &i

	fmt.Println(*p) //*represents value
	fmt.Println(&p) //&represents address

	*p = 21
	fmt.Println(*p)
}
