package notation

import "thiccgopher/game"

var LetterToPiece map[byte]game.Piece = map[byte]game.Piece{
	'p': game.BlackPawn,
	'n': game.BlackKnight,
	'b': game.BlackBishop,
	'r': game.BlackRook,
	'q': game.BlackQueen,
	'k': game.BlackKing,
	'P': game.WhitePawn,
	'N': game.WhiteKnight,
	'B': game.WhiteBishop,
	'R': game.WhiteRook,
	'Q': game.WhiteQueen,
	'K': game.WhiteKing,
}
