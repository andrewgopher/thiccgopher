package notation

import (
	"fmt"
	"thiccgopher/game"
)

func ParsePosString(posString string) *game.Pos {
	return &game.Pos{RankToInt(posString[1]), FileToInt(posString[0])}
}

func ParseMoveString(moveString string, side game.Side) *game.Move {
	move := game.Move{Start: game.Pos{RankToInt(moveString[1]), FileToInt(moveString[0])}, End: game.Pos{RankToInt(moveString[3]), FileToInt(moveString[2])}}
	if len(moveString) == 5 {
		move.Promotion = ByteToPiece[moveString[4]] //black by default
		move.Promotion -= game.Black
		move.Promotion += side
	}
	return &move
}

func MoveToUCIString(move *game.Move) string {
	result := fmt.Sprintf("%v%v%v%v", string(IntToFile(move.Start.Y)), string(IntToRank(move.Start.X)), string(IntToFile(move.End.Y)), string(IntToRank(move.End.X)))

	if move.Promotion != game.NilPiece {
		if move.Promotion&game.White > 0 {
			move.Promotion -= game.White
			move.Promotion += game.Black
		}
		result += string(PieceToByte[move.Promotion])
	}
	return result
}
