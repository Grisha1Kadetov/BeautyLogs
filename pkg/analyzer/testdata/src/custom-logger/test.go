package customlogger

import "fmt"

func CustomPrint(args ...any) {
	fmt.Println(args...)
}

func Run() {
	CustomPrint("Hello") // want "first letter should be lowercase"
	CustomPrint("привет") // want "logs should contain only English letters"
	CustomPrint("password") // want "logs should not contain sensitive data"
	CustomPrint("🤯")	// want "logs should not contain special characters"
}