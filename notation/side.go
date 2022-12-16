package notation

import "thiccgopher/game"

var ByteToSide map[byte]game.Side = map[byte]game.Side{
	'w': game.White,
	'b': game.Black,
}
