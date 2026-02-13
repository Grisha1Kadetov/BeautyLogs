package sensitive

import "fmt"

func Run() {
	fmt.Println("password") // want "logs should not contain sensitive data"
	fmt.Println("api_key=1") // want "logs should not contain sensitive data"
	fmt.Println("username: a") // want "logs should not contain sensitive data"
	fmt.Println("hello")      
}
