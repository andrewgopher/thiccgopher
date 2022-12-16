package notation

import "thiccgopher/game"

func ParseMoveString(moveString string) *game.Move {
	move := &game.Move{Start: game.Pos{RankToInt(moveString[1]), FileToInt(moveString[0])}, End: game.Pos{RankToInt(moveString[3]), FileToInt(moveString[2])}}
	if len(moveString) == 5 {
		move.Promotion = LetterToPiece[moveString[4]]
	}
	return move
}
