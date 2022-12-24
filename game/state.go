package game

type Piece = uint8

type Side = uint8

type CastleRight = uint8

//piece & White/Black to check if piece is White/Black

const (
	White Side = 64
	Black Side = 128
)

const (
	WhiteKingSide CastleRight = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

func OppSide(side Side) Side {
	if side == White {
		return Black
	} else {
		return White
	}
}

const (
	NilPiece Piece = 0

	Pawn   Piece = 1
	Knight Piece = 2
	Bishop Piece = 4
	Rook   Piece = 8
	Queen  Piece = 16
	King   Piece = 32

	WhitePawn   Piece = Pawn | White
	WhiteKnight Piece = Knight | White
	WhiteBishop Piece = Bishop | White
	WhiteRook   Piece = Rook | White
	WhiteQueen  Piece = Queen | White
	WhiteKing   Piece = King | White

	BlackPawn   Piece = Pawn | Black
	BlackKnight Piece = Knight | Black
	BlackBishop Piece = Bishop | Black
	BlackRook   Piece = Rook | Black
	BlackQueen  Piece = Queen | Black
	BlackKing   Piece = King | Black
)

type State struct {
	SideToMove      Side
	Board           [8][8]Piece
	CastleRights    CastleRight
	FiftyCount      int
	MoveCount       int
	EnPassantSquare *Pos

	PieceLists [2][]Pos
	KingPos    [2]Pos
}

func (state *State) GenPieceLists() {
	state.PieceLists = [2][]Pos{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if state.Board[i][j] != NilPiece {
				if state.Board[i][j]&White > 0 {
					state.PieceLists[0] = append(state.PieceLists[0], Pos{i, j})
				} else {
					state.PieceLists[1] = append(state.PieceLists[1], Pos{i, j})
				}
			}
			if state.Board[i][j]&King > 0 {
				if state.Board[i][j]&White > 0 {
					state.KingPos[0] = Pos{i, j}
				} else {
					state.KingPos[1] = Pos{i, j}
				}
			}
		}
	}
}

func NewState() *State {
	state := &State{}
	state.SideToMove = White
	state.Board = [8][8]Piece{{BlackRook, BlackKnight, BlackBishop, BlackQueen, BlackKing, BlackBishop, BlackKnight, BlackRook},
		{BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn},
		{NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece},
		{NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece},
		{NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece},
		{NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece, NilPiece},
		{WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn},
		{WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, WhiteBishop, WhiteKnight, WhiteRook}}
	state.CastleRights = WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide
	state.FiftyCount = 0
	state.MoveCount = 1
	state.EnPassantSquare = nil

	state.GenPieceLists()
	return state
}
