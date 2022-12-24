package game

import (
	"thiccgopher/sliceutils"
)

var (
	KingDirs   = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	KnightDirs = []Pos{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
	BishopDirs = []Pos{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	RookDirs   = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	QueenDirs  = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
)

var (
	KingSideDir  = 1
	QueenSideDir = -1
)

func SideToInd(s Side) int {
	if s == White {
		return 0
	} else {
		return 1
	}
}

func IsOnBoard(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < 8 && p.Y < 8
}

func (state *State) IsAttacked(p Pos) bool {
	//pawns
	var pawnDirection int
	var oppSide uint8
	if state.Board[p.X][p.Y]&White > 0 {
		oppSide = Black
	} else {
		oppSide = White
	}

	if oppSide == White {
		pawnDirection = -1
	} else {
		pawnDirection = 1
	}

	if p.X-pawnDirection >= 0 && p.X-pawnDirection < 8 {
		if p.Y >= 1 && state.Board[p.X-pawnDirection][p.Y-1] == oppSide|Pawn {
			return true
		}
		if p.Y <= 6 && state.Board[p.X-pawnDirection][p.Y+1] == oppSide|Pawn {
			return true
		}
	}

	for _, dir := range KingDirs {
		if IsOnBoard(Pos{p.X + dir.X, p.Y + dir.Y}) && state.Board[p.X+dir.X][p.Y+dir.Y] == oppSide|King {
			return true
		}
	}

	for _, dir := range KnightDirs {
		if IsOnBoard(Pos{p.X + dir.X, p.Y + dir.Y}) && state.Board[p.X+dir.X][p.Y+dir.Y] == oppSide|Knight {
			return true
		}
	}

	for _, dir := range BishopDirs {
		newPos := Pos{p.X + dir.X, p.Y + dir.Y}
		for IsOnBoard(newPos) {
			if state.Board[newPos.X][newPos.Y] == oppSide|Bishop || state.Board[newPos.X][newPos.Y] == oppSide|Queen {
				return true
			} else if state.Board[newPos.X][newPos.Y] != NilPiece {
				break
			}
			newPos.X += dir.X
			newPos.Y += dir.Y
		}
	}

	for _, dir := range RookDirs {
		newPos := Pos{p.X + dir.X, p.Y + dir.Y}
		for IsOnBoard(newPos) {
			if state.Board[newPos.X][newPos.Y] == oppSide|Rook || state.Board[newPos.X][newPos.Y] == oppSide|Queen {
				return true
			} else if state.Board[newPos.X][newPos.Y] != NilPiece {
				break
			}
			newPos.X += dir.X
			newPos.Y += dir.Y
		}
	}

	return false
}

func (state *State) GenPseudoMoves() []*Move { //allows the king to be in check
	moves := []*Move{}

	var currPieceList []Pos = state.PieceLists[SideToInd(state.SideToMove)]
	var pawnDirection int
	var pawnStartRank int
	var pawnPromotionRank int

	if state.SideToMove == White {
		pawnDirection = -1
		pawnStartRank = 6
		pawnPromotionRank = 0
	} else {
		pawnDirection = 1
		pawnStartRank = 1
		pawnPromotionRank = 7
	}

	currSide := state.SideToMove
	oppSide := OppSide(state.SideToMove)

	appendPawnMove := func(move *Move) { //handles promotion
		if move.End.X == pawnPromotionRank {
			moves = append(moves, &Move{move.Start, move.End, Queen | currSide})
			moves = append(moves, &Move{move.Start, move.End, Rook | currSide})
			moves = append(moves, &Move{move.Start, move.End, Bishop | currSide})
			moves = append(moves, &Move{move.Start, move.End, Knight | currSide})
		}
		moves = append(moves, move)
	}

	generateJumpMoves := func(p Pos, dirs []Pos) {
		for _, dir := range dirs {
			newPos := Pos{p.X + dir.X, p.Y + dir.Y}
			if IsOnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == NilPiece || state.Board[newPos.X][newPos.Y]&oppSide > 0) {
				moves = append(moves, &Move{Start: p, End: newPos})
			}
		}
	}
	generateSlideMoves := func(p Pos, dirs []Pos) {
		for _, dir := range dirs {
			newPos := Pos{p.X + dir.X, p.Y + dir.Y}
			for IsOnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == NilPiece || state.Board[newPos.X][newPos.Y]&oppSide > 0) {
				moves = append(moves, &Move{Start: p, End: newPos})
				newPos.X += dir.X
				newPos.Y += dir.Y
			}
		}
	}

	for _, p := range currPieceList {
		switch state.Board[p.X][p.Y] - currSide {
		case Pawn:
			if state.Board[p.X+pawnDirection][p.Y] == NilPiece {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y}})
				if p.X == pawnStartRank && state.Board[p.X+pawnDirection*2][p.Y] == NilPiece {
					appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection*2, p.Y}})
				}
			}
			if p.Y >= 1 && state.Board[p.X+pawnDirection][p.Y-1]&oppSide > 0 {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y - 1}})
			}
			if p.Y <= 6 && state.Board[p.X+pawnDirection][p.Y+1]&oppSide > 0 {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y + 1}})
			}

			if state.EnPassantSquare != nil && state.EnPassantSquare.X == p.X && (state.EnPassantSquare.Y-p.Y == 1 || state.EnPassantSquare.Y-p.Y == -1) {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y + state.EnPassantSquare.Y - p.Y}})
			}
		case King:
			generateJumpMoves(p, KingDirs)
			if !state.IsAttacked(p) {
				var KingSide, QueenSide CastleRight
				if currSide == White {
					KingSide = WhiteKingSide
					QueenSide = WhiteQueenSide
				} else {
					KingSide = BlackKingSide
					QueenSide = BlackKingSide
				}
				if state.CastleRights&KingSide > 0 && state.Board[p.X][p.Y+1] == NilPiece && state.Board[p.X][p.Y+2] == NilPiece && !state.IsAttacked(Pos{p.X, p.Y + 1}) && !state.IsAttacked(Pos{p.X, p.Y + 2}) {
					moves = append(moves, &Move{Start: p, End: Pos{p.X, p.Y + 2}})
				}
				if state.CastleRights&QueenSide > 0 && state.Board[p.X][p.Y-1] == NilPiece && state.Board[p.X][p.Y-2] == NilPiece && state.Board[p.X][p.Y-3] == NilPiece && !state.IsAttacked(Pos{p.X, p.Y - 1}) && !state.IsAttacked(Pos{p.X, p.Y - 2}) {
					moves = append(moves, &Move{Start: p, End: Pos{p.X, p.Y - 2}})
				}
			}
		case Knight:
			generateJumpMoves(p, KnightDirs)
		case Bishop:
			generateSlideMoves(p, BishopDirs)
		case Rook:
			generateSlideMoves(p, RookDirs)
		case Queen:
			generateSlideMoves(p, QueenDirs)
		default:
		}
	}
	return moves
}

func (state *State) RunMove(move *Move) (Piece, bool, int, *Pos) {
	var capturedPiece Piece
	var isEnPassant bool = false
	var oldFiftyCount int = state.FiftyCount
	var oldEnPassantSquare *Pos = state.EnPassantSquare

	//intializing info
	var pawnDirection int
	var pawnStartRank int
	var backRank int

	if state.SideToMove == White {
		pawnDirection = -1
		pawnStartRank = 6
		backRank = 7
	} else {
		pawnDirection = 1
		pawnStartRank = 1
		backRank = 0
	}

	oldStartPiece := state.Board[move.Start.X][move.Start.Y]
	oldEndPiece := state.Board[move.End.X][move.End.Y]
	currSide := state.SideToMove
	oppSide := OppSide(state.SideToMove)

	//make move

	capturedPos := move.End
	if isEnPassant {
		capturedPos = *state.EnPassantSquare
	}
	capturedPiece = state.Board[capturedPos.X][capturedPos.Y]

	state.Board[move.End.X][move.End.Y] = state.Board[move.Start.X][move.Start.Y]
	if move.Promotion != NilPiece {
		state.Board[move.End.X][move.End.Y] = move.Promotion
	}

	state.Board[move.Start.X][move.Start.Y] = NilPiece

	//en passant

	if oldStartPiece&Pawn > 0 && move.Start.Y-move.End.Y != 0 && oldEndPiece == NilPiece {
		state.Board[state.EnPassantSquare.X][state.EnPassantSquare.Y] = NilPiece
		isEnPassant = true
	}

	// piece lists

	if oldEndPiece != NilPiece {
		var oppPieceList []Pos = state.PieceLists[SideToInd(state.SideToMove)]

		for i, p := range oppPieceList {
			if p == capturedPos {
				sliceutils.RemoveByIndex(oppPieceList, i)
				break
			}
		}
	}

	for i, p := range state.PieceLists[SideToInd(currSide)] {
		if p == move.Start {
			state.PieceLists[SideToInd(currSide)][i] = move.End
			break
		}
	}

	//move counts
	if oldStartPiece&Pawn > 0 || oldEndPiece != NilPiece {
		state.FiftyCount = 0
	} else {
		state.FiftyCount++
	}

	if currSide == Black {
		state.MoveCount++
	}

	//new en passant square
	if oldStartPiece&Pawn > 0 && move.Start.X == pawnStartRank && move.End.X == pawnStartRank+pawnDirection*2 {
		state.EnPassantSquare = &move.End
	} else {
		state.EnPassantSquare = nil
	}

	//king pos
	if oldStartPiece&King > 0 {
		state.KingPos[SideToInd(currSide)] = move.End
	}

	//castle rights
	if oldStartPiece&King > 0 {
		if currSide == White {
			state.CastleRights &= ^WhiteKingSide
			state.CastleRights &= ^WhiteQueenSide
		} else {
			state.CastleRights &= ^BlackKingSide
			state.CastleRights &= ^BlackQueenSide
		}
	}
	if oldStartPiece&Rook > 0 {
		if move.Start.X == backRank {
			if move.Start.Y == 7 {
				if currSide == White {
					state.CastleRights &= ^WhiteKingSide
				} else {
					state.CastleRights &= ^BlackKingSide
				}
			} else if move.Start.Y == 0 {
				if currSide == White {
					state.CastleRights &= ^WhiteQueenSide
				} else {
					state.CastleRights &= ^BlackQueenSide
				}
			}
		}
	}

	//moving castled rook

	if oldStartPiece&King > 0 {
		if move.End.Y-move.Start.Y == 2 || move.End.Y-move.Start.Y == -2 {
			var rookPos Pos
			var newRookPos Pos
			if move.End.Y-move.Start.Y == 2 {
				rookPos = Pos{backRank, move.End.Y + 1}
				newRookPos = Pos{backRank, move.End.Y - 1}
			}
			if move.End.Y-move.Start.Y == -2 {
				rookPos = Pos{backRank, move.End.Y - 2}
				newRookPos = Pos{backRank, move.End.Y + 1}
			}
			state.Board[newRookPos.X][newRookPos.Y] = state.Board[rookPos.X][rookPos.Y]
			state.Board[rookPos.X][rookPos.Y] = NilPiece
			for i, p := range state.PieceLists[SideToInd(currSide)] {
				if p == rookPos {
					state.PieceLists[SideToInd(currSide)][i] = newRookPos
					break
				}
			}
		}
	}

	//side to move
	state.SideToMove = oppSide

	return capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantSquare
}

func (state *State) ReverseMove(move *Move, capturedPiece Piece, isEnPassant bool, oldFiftyCount int, oldEnPassantSquare *Pos) {
	oldSide := OppSide(state.SideToMove)
	var pawnDirection int
	var backRank int

	if oldSide == White {
		pawnDirection = -1
		backRank = 7
	} else {
		pawnDirection = 1
		backRank = 0
	}

	//move pieces
	oldPiece := state.Board[move.End.X][move.End.Y]
	if move.Promotion != NilPiece {
		oldPiece = Pawn & oldSide
	}

	state.Board[move.Start.X][move.Start.Y] = oldPiece

	if !isEnPassant {
		state.Board[move.End.X][move.End.Y] = capturedPiece
	} else {
		state.Board[move.End.X-pawnDirection][move.End.Y] = capturedPiece
	}

	//piece lists
	var currPieceList []Pos = state.PieceLists[SideToInd(oldSide)]
	var oppPieceListInd int = SideToInd(OppSide(oldSide))

	for i, p := range currPieceList {
		if p == move.End {
			currPieceList[i] = move.Start
		}
	}

	if capturedPiece != NilPiece {
		var capturedPiecePos Pos

		if !isEnPassant {
			capturedPiecePos = move.End
		} else {
			capturedPiecePos = Pos{move.End.X - pawnDirection, move.End.Y}
		}

		state.PieceLists[oppPieceListInd] = append(state.PieceLists[oppPieceListInd], capturedPiecePos)
	}

	//move counts

	state.FiftyCount = oldFiftyCount

	if oldSide == White {
		state.MoveCount--
	}

	//en passant square
	state.EnPassantSquare = oldEnPassantSquare

	if oldPiece&King > 0 {
		state.KingPos[SideToInd(oldSide)] = move.Start
	}

	//castling rights
	if oldPiece&King > 0 {
		if oldSide == White {
			state.CastleRights |= WhiteKingSide
			state.CastleRights |= WhiteQueenSide
		} else {
			state.CastleRights |= BlackKingSide
			state.CastleRights |= BlackQueenSide
		}
	}
	if oldPiece&Rook > 0 {
		if move.Start.X == backRank {
			if move.Start.Y == 7 {
				if oldSide == White {
					state.CastleRights |= WhiteKingSide
				} else {
					state.CastleRights |= BlackKingSide
				}
			} else if move.Start.Y == 0 {
				if oldSide == White {
					state.CastleRights |= WhiteQueenSide
				} else {
					state.CastleRights |= BlackQueenSide
				}
			}
		}
	}

	//reverse castling

	if oldPiece&King > 0 {
		if move.End.Y-move.Start.Y == 2 || move.End.Y-move.Start.Y == -2 {
			var rookPos Pos
			var newRookPos Pos
			if move.End.Y-move.Start.Y == 2 {
				rookPos = Pos{backRank, move.End.Y + 1}
				newRookPos = Pos{backRank, move.End.Y - 1}
			}
			if move.End.Y-move.Start.Y == -2 {
				rookPos = Pos{backRank, move.End.Y - 2}
				newRookPos = Pos{backRank, move.End.Y + 1}
			}

			state.Board[rookPos.X][rookPos.Y] = state.Board[newRookPos.X][newRookPos.Y]
			state.Board[newRookPos.X][newRookPos.Y] = NilPiece

			for i, p := range state.PieceLists[SideToInd(oldSide)] {
				if p == newRookPos {
					state.PieceLists[SideToInd(oldPiece)][i] = rookPos
					break
				}
			}
		}
	}

	//side to move
	state.SideToMove = oldSide
}

func (state *State) GenMoves() []*Move { //all valid moves
	pseudoMoves := state.GenPseudoMoves()
	validMoves := []*Move{}
	currSide := state.SideToMove

	for _, pseudoMove := range pseudoMoves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantSquare := state.RunMove(pseudoMove)

		if !state.IsAttacked(state.KingPos[SideToInd(currSide)]) {
			validMoves = append(validMoves, pseudoMove)
		}

		state.ReverseMove(pseudoMove, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantSquare)
	}

	return validMoves
}
