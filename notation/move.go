package notation

import "thiccgopher/game"

func ParseSquareString(squareString string) *game.Pos {
	return &game.Pos{RankToInt(squareString[1]), FileToInt(squareString[0])}
}

func ParseMoveString(moveString string, side game.Side) *game.Move {
	move := &game.Move{Start: &game.Pos{RankToInt(moveString[1]), FileToInt(moveString[0])}, End: &game.Pos{RankToInt(moveString[3]), FileToInt(moveString[2])}}
	if len(moveString) == 5 {
		move.Promotion = ByteToPiece[moveString[4]] //black by default
		move.Promotion -= game.Black
		move.Promotion += side
	}
	return move
}
