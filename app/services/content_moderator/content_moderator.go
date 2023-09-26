package helpers

import "fmt"

type Policy struct {
	ID      int
	Name    string
	Content string
}

func main() {
	openaiConnect()
}

func openaiConnect() {
	fmt.Println("I'll do this tomorrow")
	//TODO
	//fine tune  model
	//feed from docs/ + promt
	//test results
}
