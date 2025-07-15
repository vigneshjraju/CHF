package main

import "fmt"

func main() {

	phonebook := map[string]string{"Tony": "5764784"}
	phonebook["Tommy"] = "12345"
	phonebook["Danny"] = "68521"

	fmt.Println(phonebook)
	fmt.Println(phonebook["Tommy"])

}
