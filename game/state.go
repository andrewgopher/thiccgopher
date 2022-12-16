package game

type Piece int8

type Side int8

//piece & White/Black to check if piece is White/Blac

var White Side = 8
var Black Side = 16

var (
	NilPiece Piece = 0

	WhitePawn   Piece = 9
	WhiteKnight Piece = 10
	WhiteBishop Piece = 11
	WhiteRook   Piece = 12
	WhiteQueen  Piece = 13
	WhiteKing   Piece = 14

	BlackPawn   Piece = 17
	BlackKnight Piece = 18
	BlackBishop Piece = 19
	BlackRook   Piece = 20
	BlackQueen  Piece = 21
	BlackKing   Piece = 22
)

type State struct {
	SideToMove   Side
	Board        [8][8]Piece
	CastleRights [2][2]bool
	FiftyCount   int
	MoveCount    int
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
	state.CastleRights = [2][2]bool{{true, true}, {true, true}}
	state.FiftyCount = 0
	state.MoveCount = 1
	return state
}
