package main

import "fmt"

func main() {
	fmt.Println("Welcome to Go programming language.")
	var name string
	fmt.Print("Please enter your Name:")
	fmt.Scan(&name)
	fmt.Printf("Hello %s goodmorning \n", name)
}
