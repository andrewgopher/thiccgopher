package notation

import (
	"strconv"
	"strings"
	"thiccgopher/game"
)

func ParseFenString(fen string) *game.State {
	state := &game.State{}
	args := strings.Split(fen, " ")
	ranks := strings.Split(args[0], "/")

	for i, rank := range ranks {
		currFile := 0
		for _, c := range rank {
			if byte(c) >= '1' && byte(c) <= '8' {
				currFile += int(c-'1') + 1
			} else {
				state.Board[i][currFile] = ByteToPiece[byte(c)]
				currFile += 1
			}
		}
	}

	state.SideToMove = ByteToSide[args[1][0]]

	state.CastleRights = [2][2]bool{{strings.Contains(args[2], "K"), strings.Contains(args[2], "Q")}, {strings.Contains(args[3], "k"), strings.Contains(args[3], "q")}}
	state.EnPassantSquare = nil
	if args[3] != "-" {
		state.EnPassantSquare = ParseSquareString(args[3])
	}

	state.FiftyCount, _ = strconv.Atoi(args[4])

	state.MoveCount, _ = strconv.Atoi(args[5])

	return state
}
