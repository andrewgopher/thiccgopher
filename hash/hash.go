package hash

import (
	"math/rand"
	"thiccgopher/game"
)

var pieceTable [2][6][8][8]uint64
var castleRightTable [4]uint64
var enPassantTable [8][8]uint64
var sideToMoveTable [2]uint64
var hasInit bool

func genBitTables(currPieceTable *[2][6][8][8]uint64, currCastleRightTable *[4]uint64, currEnPassantTable *[8][8]uint64, currSideToMoveTable *[2]uint64) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 6; j++ {
			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					currPieceTable[i][j][x][y] = rand.Uint64()
				}
			}
		}
	}

	for i := 0; i < 4; i++ {
		currCastleRightTable[i] = rand.Uint64()
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			currEnPassantTable[i][j] = rand.Uint64()
		}
	}
	for i := 0; i < 2; i++ {
		currSideToMoveTable[i] = rand.Uint64()
	}
}

func initBitTables() {
	hasInit = true
	genBitTables(&pieceTable, &castleRightTable, &enPassantTable, &sideToMoveTable)
}

func Hash(state *game.State) uint64 { //TODO: return two hashes for less collisions
	if !hasInit {
		initBitTables()
	}
	var result uint64 = 0

	//pieces
	for i := 0; i < 2; i++ {
		for _, p := range state.PieceLists[i] {
			result ^= pieceTable[i][game.PieceToInd[game.PieceOnly(state.Board[p.X][p.Y])]][p.X][p.Y]
		}
	}

	//castle rights
	for _, c := range game.CastleRights {
		if state.CastleRights&c > 0 {
			result ^= castleRightTable[game.CastleRightToInd[c]]
		}
	}

	//en passant
	if state.EnPassantPos != nil {
		result ^= enPassantTable[state.EnPassantPos.X][state.EnPassantPos.Y]
	}

	//side to move
	result ^= castleRightTable[game.SideToInd[state.SideToMove]]

	return result
}
