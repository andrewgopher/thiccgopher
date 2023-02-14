package engine

import (
	"math"
	"thiccgopher/game"
)

const CheckmateEval = 1000000
const CheckmateEvalSplit = 900000 //if absolute value is above this, then it is a checkmate eval. otherwise, normal or stalemate eval
const BigEval = math.MaxInt / 2

var PieceValues map[game.Piece]int = map[game.Piece]int{
	game.Pawn:   100,
	game.Bishop: 300,
	game.Knight: 300,
	game.Rook:   500,
	game.Queen:  900,
}

func GetKingStatus(state *game.State, side game.Side) (bool, bool) {
	kingPos := state.KingPos[game.SideToInd[side]]
	oppSide := game.OppSide(game.PieceSide(state.Board[kingPos.X][kingPos.Y]))
	isInCheck := state.IsAttacked(kingPos, oppSide)
	isSurrounded /*by check*/ := true
	for _, dir := range game.KingDirs {
		newPos := game.Pos{kingPos.X + dir.X, kingPos.Y + dir.Y}
		if game.IsOnBoard(newPos) && state.Board[newPos.X][newPos.Y]&side == 0 { //square is empty or contains an opponent's piece
			if !state.IsAttacked(newPos, oppSide) {
				isSurrounded = false
				break
			}
		}
	}
	return isInCheck, isSurrounded
}

func Eval(state *game.State) (int, bool) { //int is eval, bool is if decisive
	side := state.SideToMove
	numMoves := len(state.GenMoves())
	if state.FiftyCount == 50 {
		return 0, true
	}

	if len(state.PieceLists[0]) == 1 && len(state.PieceLists[1]) == 1 {
		return 0, true
	}

	if numMoves == 0 {
		isOppInCheck := state.IsAttacked(state.KingPos[game.SideToInd[game.OppSide(side)]], side)

		if !isOppInCheck { //stalemate
			return 0, true
		} else { //checkmate
			return CheckmateEval, true
		}
	}

	materialEval := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pieceValue := PieceValues[game.PieceOnly(state.Board[i][j])]
			if state.Board[i][j]&side > 0 {
				materialEval += pieceValue
			} else {
				materialEval -= pieceValue
			}
		}
	}

	pieceMapEval := 0
	for _, side := range game.Sides {
		for _, p := range state.PieceLists[game.SideToInd[side]] {
			rank := p.X
			if side == game.Black {
				rank = 7 - p.X
			}
			posEval := pieceMaps[game.PieceOnly(state.Board[p.X][p.Y])][rank][p.Y]
			if side == state.SideToMove {
				pieceMapEval += posEval
			} else {
				pieceMapEval -= posEval
			}
		}
	}

	return materialEval + pieceMapEval/2, false
}
