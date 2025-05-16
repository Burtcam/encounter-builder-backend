package main

import "fmt"
import "github.com/Burtcam/encounter-builder-backend/utils"

// struct encounter {
// 	difficulty string
// 	pSize      int
// 	level      int
// }

func main() {
	difficulty := getXpBudget("Trivial", 4, 1)
	fmt.Println(difficulty)
}
