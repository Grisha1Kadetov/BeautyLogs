package lowercase

import "fmt"

func Run() {
	fmt.Println("Hello") // want "first letter should be lowercase"
	fmt.Println("hello")
	fmt.Println("Привет") // want "first letter should be lowercase"
	fmt.Println("привет")
}