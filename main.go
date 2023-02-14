package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"thiccgopher/boolwrapper"
	"thiccgopher/deepcopy"
	"thiccgopher/engine"
	"thiccgopher/game"
	"thiccgopher/notation"
	"time"
)

const (
	name   = "thiccgopher"
	author = "andrewgopher"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var logFileName string
var logFile io.Writer

func GetInputLn() string {
	if scanner.Scan() {
		line := scanner.Text()
		if logFile != nil {
			fmt.Fprintln(logFile, line)
		}
		return line
	}
	return ""
}

func setPosition(options []string) *game.State {
	var state *game.State
	if options[0] == "startpos" {
		state = game.NewState()
	} else {
		state = notation.ParseFenString(options[0])
	}
	if len(options) > 2 {
		for i := 2; i < len(options); i++ {
			state.RunMove(notation.ParseMoveString(options[i], state.SideToMove))
		}
	}
	return state
}

func Perft(state *game.State, depth int) uint64 {
	if depth == 0 {
		return 1
	}
	moves := state.GenMoves()
	if depth == 1 {
		return uint64(len(moves))
	} else {
		var totalMoves uint64 = 0
		for _, m := range moves {
			capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
			currMoves := Perft(state, depth-1)
			totalMoves += currMoves
			state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		}
		return totalMoves
	}
}

func StopAfterDuration(timeLimit time.Duration, isSearching *boolwrapper.BoolWrapper) {
	time.Sleep(timeLimit)
	isSearching.Val = false
}

func ProcessSearch(state *game.State, timeLimit time.Duration, moveChan chan *game.Move, isSearching *boolwrapper.BoolWrapper) {
	copiedStateIface, _ := deepcopy.Anything(state)
	go engine.IterativeDeepening(copiedStateIface.(*game.State), moveChan, isSearching)
	go StopAfterDuration(timeLimit, isSearching)
	for isSearching.Val {

	}
	for len(moveChan) > 1 {
		<-moveChan
	}
	bestMove := <-moveChan
	fmt.Println("bestmove", notation.MoveToUCIString(bestMove))
	isSearching.Val = false
}

func main() {
	var state *game.State
	var moveChan chan *game.Move = make(chan *game.Move, 1000)
	isSearching := &boolwrapper.BoolWrapper{Val: false}

	flag.StringVar(&logFileName, "logFile", "", "")
	flag.Parse()

	if logFileName != "" {
		logFile, _ = os.Create(logFileName)
	}

	rand.Seed(time.Now().UnixNano())
	var input string

	//check if uci
	input = GetInputLn()
	if input != "uci" {
		fmt.Println("Only UCI mode supported")
		os.Exit(1)
	}

	//identify

	fmt.Printf("id name %v\nid author %v\n", name, author)

	//engine options

	//engine options done

	fmt.Println("uciok")

	for {
		input = GetInputLn()
		tokens := strings.Split(input, " ")
		command := tokens[0]
		if isSearching.Val {
			switch command {
			case "quit":
				os.Exit(0)
			case "stop":
				isSearching.Val = false
			}
			continue
		}
		switch command {
		case "quit":
			os.Exit(0)
		case "position":
			options := []string{}
			if tokens[1] == "startpos" {
				options = tokens[1:]
			} else if tokens[1] == "fen" {
				options = append(options, "")
				for i := 2; i <= 7; i++ {
					options[0] += tokens[i] + " "
				}
				if len(tokens) > 8 {
					options = append(options, tokens[8:]...)
				}
			}
			state = setPosition(options)
		case "go":
			if len(tokens) > 1 && tokens[1] == "perft" {
				depth, _ := strconv.Atoi(tokens[2])
				moves := state.GenMoves()
				var totalPerft uint64 = 0
				for _, m := range moves {
					capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
					currPerft := Perft(state, depth-1)
					fmt.Println(notation.MoveToUCIString(m), currPerft)
					state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
					totalPerft += currPerft
				}
				fmt.Println(totalPerft)
			} else {
				var options map[string]string = map[string]string{}

				for i := 0; i < (len(tokens)-1)/2; i++ {
					options[tokens[1+2*i]] = tokens[2+2*i]
				}

				var whiteTime int
				var blackTime int
				var whiteInc int
				var blackInc int
				for k, v := range options {
					if k == "wtime" {
						whiteTime, _ = strconv.Atoi(v)
					} else if k == "btime" {
						blackTime, _ = strconv.Atoi(v)
					} else if k == "movetime" {
						whiteTime, _ = strconv.Atoi(v)
						blackTime, _ = strconv.Atoi(v)
					} else if k == "winc" {
						whiteInc, _ = strconv.Atoi(v)
					} else if k == "binc" {
						blackInc, _ = strconv.Atoi(v)
					}
				}

				var timeLimit time.Duration
				var timeLeft time.Duration
				var timeInc time.Duration

				if state.SideToMove == game.White {
					timeLeft = time.Duration(whiteTime) * time.Millisecond
					timeInc = time.Duration(whiteInc) * time.Millisecond
				} else {
					timeLeft = time.Duration(blackTime) * time.Millisecond
					timeInc = time.Duration(blackInc) * time.Millisecond
				}

				if timeLeft == 0 {
					timeLimit = math.MaxInt64
				} else {
					timeLimit = timeLeft/40 + timeInc
					if timeLimit > timeLeft*9/10 {
						timeLimit = timeLeft * 9 / 10
					}
				}
				isSearching.Val = true
				go ProcessSearch(state, timeLimit, moveChan, isSearching)
			}
		case "ucinewgame":
		case "isready":
			fmt.Println("readyok")
		default:
			fmt.Println("Unsupported command")
			os.Exit(1)
		}
	}
}
