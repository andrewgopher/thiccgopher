package engine

import (
	"sort"
	"thiccgopher/game"
	"thiccgopher/sliceutils"
	"time"
)

func Search(state *game.State, timeLimit time.Duration) *game.Move {
	currDepth := 2
	moves := state.GenMoves()
	startTime := time.Now()
	var bestMove *game.Move

	var prevSearchDuration time.Duration
	var prevPrevSearchDuration time.Duration

	for time.Since(startTime) < timeLimit && (prevPrevSearchDuration == 0 || prevSearchDuration == 0 || timeLimit-time.Since(startTime) >= prevSearchDuration*prevSearchDuration/prevPrevSearchDuration) {
		currSearchStartTime := time.Now()

		bestMove = nil
		var bestMoveInd int
		bestEval := -BigEval
		for i, m := range moves {
			capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
			currScore := -Minimax(state, currDepth-1, -BigEval, -bestEval)
			if currScore > bestEval {
				bestEval = currScore
				bestMove = m
				bestMoveInd = i
			}
			state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		}

		moves = sliceutils.RemoveByIndex(moves, bestMoveInd)
		moves = append([]*game.Move{bestMove}, moves...)

		prevPrevSearchDuration = prevSearchDuration
		prevSearchDuration = time.Since(currSearchStartTime)

		currDepth++
	}
	return bestMove
}

func Minimax(state *game.State, depth int, alpha, beta int) int {
	currSide := state.SideToMove
	moves := state.GenMoves()

	bestEval := -BigEval
	if len(moves) == 0 {
		if state.IsAttacked(state.KingPos[game.SideToInd(currSide)], game.OppSide(currSide)) {
			bestEval = -CheckmateEval
		} else {
			bestEval = 0
		}
	}
	sort.Slice(moves, func(a int, b int) bool {
		scoreA := PieceValues[game.PieceOnly(state.Board[moves[a].End.X][moves[a].End.Y])] - PieceValues[game.PieceOnly(state.Board[moves[a].Start.X][moves[a].Start.Y])]/100
		scoreB := PieceValues[game.PieceOnly(state.Board[moves[b].End.X][moves[b].End.Y])] - PieceValues[game.PieceOnly(state.Board[moves[b].Start.X][moves[b].Start.Y])]/100
		return scoreA > scoreB
	})
	for _, m := range moves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
		var currOppEval int
		if depth == 1 {
			currOppEval, _ = Eval(state) //evaluates in the pov of opponent
		} else {
			currOppEval = Minimax(state, depth-1, -beta, -alpha) //evaluates in the pov of opponent
		}
		currEval := -currOppEval
		if currEval > bestEval {
			bestEval = currEval
		}
		if bestEval > alpha {
			alpha = bestEval
		}
		state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		if alpha >= beta {
			break
		}
	}
	if bestEval > CheckmateEvalSplit {
		bestEval -= 1
	}
	if bestEval < -CheckmateEvalSplit {
		bestEval += 1
	}
	return bestEval
}
