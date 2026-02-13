package special

import "fmt"

func Run() {
	fmt.Println("hello!") // want "logs should not contain special characters"
	fmt.Println(".") // want "logs should not contain special characters"
	fmt.Println("🤯") // want "logs should not contain special characters"
	fmt.Println("#") //ignore
}
