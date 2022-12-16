package main

import (
	"fmt"
)

var name = "thiccgopher"
var author = "andrewgopher"

func main() {
	var input string

	//check if uci
	fmt.Scanln(&input)
	if input != "uci" {
		panic("Only UCI mode supported")
	}

	//identify

	fmt.Printf("id name %v\nid author %v\n", name, author)

	//engine options

	//engine options done

	fmt.Println("uciok")

	//GUI options
	fmt.Scanln(&input)

	if input == "quit" {
		return
	} else if input == "isready" {

	} else {
		panic("Unsupported command")
	}

	//GUI options done
	fmt.Println("readyok")
}
