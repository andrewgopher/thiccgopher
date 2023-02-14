package engine

import (
	"sort"
	"thiccgopher/boolwrapper"
	"thiccgopher/game"
	"thiccgopher/hash"
	"thiccgopher/sliceutils"
)

type pvEntry struct {
	move  *game.Move
	depth int
}

var PVS map[uint64]*pvEntry = make(map[uint64]*pvEntry)

func IterativeDeepening(state *game.State, moveChan chan *game.Move, isSearching *boolwrapper.BoolWrapper) {
	currDepth := 2
	moves := state.GenMoves()
	var bestMove *game.Move
	for {
		_, bestMove = Minimax(state, currDepth, -BigEval, BigEval, isSearching)
		if !isSearching.Val {
			return
		}
		var bestMoveInd int
		for i := range moves {
			if moves[i] == bestMove {
				bestMoveInd = i
			}
		}

		moves = sliceutils.RemoveByIndex(moves, bestMoveInd)
		moves = append([]*game.Move{bestMove}, moves...)
		moveChan <- bestMove
		currDepth++
	}
}

func Minimax(state *game.State, depth int, alpha, beta int, isSearching *boolwrapper.BoolWrapper) (int, *game.Move) {
	currSide := state.SideToMove
	currHash := hash.Hash(state)
	moves := []*game.Move{}

	pv, hasStoredPV := PVS[currHash]
	if hasStoredPV {
		moves = append(moves, pv.move)
	}
	moves = append(moves, state.GenMoves()...)

	bestEval := -BigEval
	var bestMove *game.Move = nil
	if len(moves) == 0 {
		if state.IsAttacked(state.KingPos[game.SideToInd[currSide]], game.OppSide(currSide)) {
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
			currOppEval, _ = Minimax(state, depth-1, -beta, -alpha, isSearching) //evaluates in the pov of opponent
		}
		if !isSearching.Val {
			return 0, nil
		}
		currEval := -currOppEval
		if currEval > bestEval {
			bestEval = currEval
			bestMove = m
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
	if hasStoredPV && depth > pv.depth {
		PVS[currHash] = &pvEntry{bestMove, depth}
	}
	return bestEval, bestMove
}
