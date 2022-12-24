package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"thiccgopher/game"
	"thiccgopher/notation"
)

const (
	name   = "thiccgopher"
	author = "andrewgopher"
)

var engineChan chan game.Move = make(chan game.Move) //engine-chan uwu
var currState *game.State

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var f, _ = os.Create("/tmp/thicclog.txt")

func GetInputLn() string {
	if scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(f, line)
		return line
	}
	return ""
}

func setPosition(options []string) {
	if options[0] == "startpos" {
		currState = game.NewState()
	} else {
		currState = notation.ParseFenString(options[0])
	}
	if len(options) > 2 {
		for i := 2; i < len(options); i++ {
			currState.RunMove(notation.ParseMoveString(options[i], currState.SideToMove))
		}
	}
}

func main() {
	var input string

	//check if uci
	input = GetInputLn()
	if input != "uci" {
		panic("Only UCI mode supported")
	}

	//identify

	fmt.Printf("id name %v\nid author %v\n", name, author)

	//engine options

	//engine options done

	fmt.Println("uciok")

	//GUI options
	input = GetInputLn()

	if input == "quit" || input == "stop" {
		return
	} else if input == "isready" {

	} else if input == "ucinewgame" {
	} else {
		panic("Unsupported command")
	}

	//GUI options done
	fmt.Println("readyok")

	for {
		input = GetInputLn()
		tokens := strings.Split(input, " ")
		command := tokens[0]
		switch command {
		case "quit", "stop":
			os.Exit(0)
		case "position":
			options := []string{}
			if tokens[1] == "startpos" {
				options = tokens[1:]
			} else {
				options = append(options, "")
				for i := 1; i <= 6; i++ {
					options[0] += tokens[i] + " "
				}
				if len(tokens) > 7 {
					options = append(options, tokens[7:]...)
				}
			}
			setPosition(options)
		case "go":
			fmt.Println("bestmove", notation.MoveToFenString(currState.GenMoves()[0]))
		case "ucinewgame":
		case "isready":

		default:
			panic("Unsupported command")
		}
	}
}
