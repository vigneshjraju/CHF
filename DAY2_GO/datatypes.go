package main

import "fmt"

func main() {
	var number1 int
	var number2 float64
	var complex complex64
	var name string
	var IsTrue bool

	fmt.Println(number1, number2, complex, name, IsTrue)

	var age int                  //declaration without initialisation
	var user = "Tommy"           //declaration with initialisation
	email := "vignesh@gmail.com" //shorthand for declaration with initialisation
	fmt.Println(age, user, email)

	const distance = 25
	fmt.Printf("Type of distance: %T \n", distance)
	fmt.Printf("Type of distance: %v \n", distance)

}
