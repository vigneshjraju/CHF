package main

import (
	"fmt"
	"time"
)

func main() {

	day := 2

	switch day {

	case 1:
		fmt.Println("it is monday")

	case 2:
		fmt.Println("it is tuesday")

	case 3:
		fmt.Println("it is wednesday")

	default:
		fmt.Println("weekend")

	}

	now := time.Now()

	switch {

	case now.Hour() < 12:
		fmt.Println("Goodmorning")

	case now.Hour() < 16:
		fmt.Println("Goodafternoon")

	default:
		fmt.Println("Goodevening")

	}

}
