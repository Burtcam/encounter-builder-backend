package main

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

func main() {
	content, err := os.ReadFile("example-beast.json")
	if err != nil {
		fmt.Println("Fucked it up")
	}
	fmt.Println(string(content))
	fmt.Println(gjson.Get(string(content), "name"))
	fmt.Println(gjson.Get(string(content), "system.abilities.cha.mod"))

}
