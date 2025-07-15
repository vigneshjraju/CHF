package main

import "fmt"

type User struct {
	Name string
	Age int
	Email string
}

func main(){
	user1 := User{"Vignesh",27,"vigneshjraju@gmail.com"}
	user2 := User{"Bal",28,"sarathkbalan@gmail.com"}

	fmt.Println(user1,user2)
	fmt.Println(user2.Email)
	fmt.Println(user2.Age)

	

}