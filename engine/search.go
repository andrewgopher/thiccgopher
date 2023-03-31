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

var PVS map[uint64]map[uint64]*pvEntry = make(map[uint64]map[uint64]*pvEntry)
var bitTables1 = hash.NewBitTables()
var bitTables2 = hash.NewBitTables()

func CaptureSearch(state *game.State, alpha, beta int) int {
	standPat, isDecisive := Eval(state)
	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}
	if isDecisive {
		return standPat
	}
	captures := state.GenPseudoMoves()
	selfSide := state.SideToMove
	for _, capture := range captures {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(capture)
		if state.IsAttacked(state.KingPos[game.SideToInd[selfSide]], game.OppSide(selfSide)) || capturedPiece == game.NilPiece {
			state.ReverseMove(capture, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
			continue
		}
		eval := -CaptureSearch(state, -beta, -alpha)
		state.ReverseMove(capture, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		if eval >= beta {
			return beta
		}
		if eval > alpha {
			alpha = eval
		}
	}
	return alpha
}

func IterativeDeepening(state *game.State, moveChan chan *game.Move, isSearching *boolwrapper.BoolWrapper) {
	currDepth := 2
	moves := state.GenMoves()
	var bestMove *game.Move
	for {
		// currDepthStartTime := time.Now()
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
		// fmt.Println(currDepth, time.Since(currDepthStartTime))
		currDepth++
	}
}

func Minimax(state *game.State, depth int, alpha, beta int, isSearching *boolwrapper.BoolWrapper) (int, *game.Move) {
	currSide := state.SideToMove
	currHash1 := hash.Hash(state, bitTables1)
	currHash2 := hash.Hash(state, bitTables2)
	moves := state.GenPseudoMoves()

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
	pv, hasStoredPV := PVS[currHash1]
	if hasStoredPV {
		moves = append([]*game.Move{pv[currHash2].move}, moves...)
	}
	selfSide := state.SideToMove
	for _, m := range moves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
		if state.IsAttacked(state.KingPos[game.SideToInd[selfSide]], game.OppSide(selfSide)) {
			state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
			continue
		}
		var currOppEval int
		if depth == 1 {
			currOppEval = CaptureSearch(state, -beta, -alpha) //evaluates in the pov of opponent
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
	if bestMove != nil && (!hasStoredPV || depth > pv[currHash2].depth) {
		if hasStoredPV {
			PVS[currHash1][currHash2] = &pvEntry{bestMove, depth}
		} else {
			PVS[currHash1] = make(map[uint64]*pvEntry)
			PVS[currHash1][currHash2] = &pvEntry{bestMove, depth}
		}
	}
	return bestEval, bestMove
}
