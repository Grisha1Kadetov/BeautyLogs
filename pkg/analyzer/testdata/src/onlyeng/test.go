package onlyeng

import "fmt"

func Run() {
	fmt.Println("hello")
	fmt.Println("привет") // want "logs should contain only English letters"
	fmt.Println("helloПривет") // want "logs should contain only English letters"
}
